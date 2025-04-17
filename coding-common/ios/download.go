package ios

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type writeCounter struct {
	Total   uint64
	Path    string
	written uint64
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.written += uint64(n)
	percent := float64(wc.written) / float64(wc.Total) * 100
	fmt.Printf("\r下载进度 %s: %.2f%%", wc.Path, percent)
	return n, nil
}

type sectionWriter struct {
	file   *os.File
	offset int64
}

func (w *sectionWriter) Write(p []byte) (n int, err error) {
	n, err = w.file.WriteAt(p, w.offset)
	w.offset += int64(n)
	return
}

// 流式下载
func DownloadFile(url, filepath string) error {
	// 发起GET请求
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// 创建文件
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将响应体直接写入文件
	_, err = io.Copy(file, resp.Body)
	return err
}

// 带进度条
func DownloadWithProgress(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return err
	}

	tmpFilepath := filepath + ".tmp"
	file, err := os.Create(tmpFilepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, io.TeeReader(resp.Body, &writeCounter{
		Total: uint64(size),
		Path:  filepath,
	}))
	if err != nil {
		return err
	}
	return os.Rename(tmpFilepath, filepath)
}

// 支持断点续传
func DownloadWithResume(url, filepath string) error {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return err
	}

	// 设置Range头
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if fi.Size() > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", fi.Size()))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 如果是部分内容响应(206)或完整响应(200)
	if resp.StatusCode != http.StatusPartialContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %s", resp.Status)
	}

	// 移动到文件末尾继续写入
	_, err = file.Seek(fi.Size(), io.SeekStart)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	return err
}

// 多线程/分段下载
func DownloadConcurrent(url, filepath string, workers int) error {
	// 获取文件大小
	resp, err := http.Head(url)
	if err != nil {
		return err
	}
	size, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return err
	}

	// 创建文件
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 设置文件大小
	if err := file.Truncate(size); err != nil {
		return err
	}

	// 计算每个worker负责的字节范围
	chunkSize := size / int64(workers)
	errChan := make(chan error, workers)
	var wg sync.WaitGroup

	for i := range workers {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			start := int64(i) * chunkSize
			end := start + chunkSize - 1
			if i == workers-1 {
				end = size - 1 // 最后一个worker处理剩余部分
			}

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				errChan <- err
				return
			}

			req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				errChan <- err
				return
			}
			defer resp.Body.Close()

			// 写入文件的指定位置
			writer := &sectionWriter{
				file:   file,
				offset: start,
			}

			if _, err := io.Copy(writer, resp.Body); err != nil {
				errChan <- err
				return
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func RobustDownload(url, filepath string, maxRetries int) error {
	tmpPath := filepath + ".tmp"

	var err error
	for i := range maxRetries {
		err = DownloadFile(url, tmpPath)
		if err == nil {
			// 下载成功，重命名文件
			return os.Rename(tmpPath, filepath)
		}
		time.Sleep(time.Second * time.Duration(i+1)) // 指数退避
	}

	// 所有重试都失败后清理临时文件
	os.Remove(tmpPath)
	return fmt.Errorf("after %d attempts, last error: %v", maxRetries, err)
}
