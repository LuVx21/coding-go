package ios

import (
	"bufio"
	"io"
	"os"
	"sync"

	"golang.org/x/sync/errgroup"
)

// WriteOptions 写文件的参数
type WriteOptions struct {
	BufferSize  int         // 缓冲区大小(字节), 默认 4KB
	Concurrency int         // 并发写入协程数(仅限 WriteChunks 使用)
	Permission  os.FileMode // 文件权限, 默认 0644
}

// WriteFile 高效写入文件(支持大文件分块) 内存占用: 高
func WriteFile(filePath string, data []byte, opts *WriteOptions) error {
	if opts == nil {
		opts = &WriteOptions{
			BufferSize: 4096,
			Permission: 0644,
		}
	}

	// 原子性写入(通过临时文件+重命名)
	tmpFile := filePath + ".tmp"
	f, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, opts.Permission)
	if err != nil {
		return err
	}
	defer f.Close()

	// 缓冲写入
	writer := bufio.NewWriterSize(f, opts.BufferSize)
	if _, err := writer.Write(data); err != nil {
		_ = os.Remove(tmpFile)
		return err
	}

	// 确保缓冲数据刷盘
	if err := writer.Flush(); err != nil {
		_ = os.Remove(tmpFile)
		return err
	}

	// 原子替换原文件
	return os.Rename(tmpFile, filePath)
}

// WriteStream 流式写入(适合网络传输或大文件) 内存占用: 低
func WriteStream(filePath string, reader io.Reader, opts *WriteOptions) error {
	if opts == nil {
		opts = &WriteOptions{
			BufferSize: 4096,
			Permission: 0644,
		}
	}

	tmpFile := filePath + ".tmp"
	f, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, opts.Permission)
	if err != nil {
		return err
	}
	defer f.Close()

	// 使用 bufio 加速写入
	bufWriter := bufio.NewWriterSize(f, opts.BufferSize)
	if _, err := io.Copy(bufWriter, reader); err != nil {
		_ = os.Remove(tmpFile)
		return err
	}

	if err := bufWriter.Flush(); err != nil {
		_ = os.Remove(tmpFile)
		return err
	}

	return os.Rename(tmpFile, filePath)
}

// WriteChunks 并发分块写入(适合超大文件) 内存占用: 中
func WriteChunks(filePath string, chunks <-chan []byte, opts *WriteOptions) error {
	if opts == nil {
		opts = &WriteOptions{
			BufferSize:  4096,
			Concurrency: 4,
			Permission:  0644,
		}
	}

	tmpFile := filePath + ".tmp"
	f, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, opts.Permission)
	if err != nil {
		return err
	}
	defer f.Close()

	var (
		mu     sync.Mutex
		writer = bufio.NewWriterSize(f, opts.BufferSize)
		g      errgroup.Group
	)

	for range opts.Concurrency {
		g.Go(func() error {
			for chunk := range chunks {
				mu.Lock()
				_, err := writer.Write(chunk)
				mu.Unlock()
				if err != nil {
					return err
				}
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		_ = os.Remove(tmpFile)
		return err
	}

	if err := writer.Flush(); err != nil {
		_ = os.Remove(tmpFile)
		return err
	}

	return os.Rename(tmpFile, filePath)
}
