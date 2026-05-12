package main

import (
	"strings"
)

const (
	default_registry = "registry-1.docker.io"
	default_repo     = "library"
	default_tag      = "latest"
)

type ImageInfo struct {
	Full             string // 原始完整镜像名
	Registry         string // registry 地址，如 docker.io, harbor.example.com:5000
	Namespace        string // 命名空间/组织，如 library, myorg
	Image            string // 镜像名，如 nginx, mysql
	Tag              string // 标签，如 latest, v1.0.0
	Digest           string // 摘要（如果有），如 sha256:abc123...
	Os, Architecture string
	Size             int64

	NewerTag string
}

func pickTag(full string) string {
	parts := strings.Split(full, ":")
	if len(parts) <= 1 {
		return default_tag
	}
	tag := parts[len(parts)-1]
	if strings.Contains(tag, "/") && tag != "" {
		return default_tag
	}
	return tag
}
func pickDigest(s string) string { return s[strings.Index(s, "@")+1:] }

func parseImage(full string) ImageInfo {
	info := ImageInfo{Registry: default_registry, Namespace: default_repo, Tag: default_tag, Full: full}
	img := full

	parts := strings.Split(img, ":")
	if len(parts) > 1 {
		last := parts[len(parts)-1]
		if !strings.Contains(last, "/") && last != "" {
			info.Tag = last
			img = strings.Join(parts[:len(parts)-1], ":")
		}
	}

	parts = strings.Split(img, "/")
	if len(parts) == 1 {
		info.Image = parts[0]
	} else if len(parts) == 2 {
		// 两种可能: registry/image 或 namespace/image
		if strings.Contains(parts[0], ".") || strings.Contains(parts[0], ":") || parts[0] == "localhost" {
			// registry/image
			info.Registry = parts[0]
			info.Image = parts[1]
		} else {
			// namespace/image
			info.Namespace = parts[0]
			info.Image = parts[1]
		}
	} else if len(parts) >= 3 {
		// registry/namespace/image 或 registry/ns1/ns2/image
		info.Registry = parts[0]
		info.Namespace = strings.Join(parts[1:len(parts)-1], "/")
		info.Image = parts[len(parts)-1]
	}
	return info
}
func extractValue(header, key string) string {
	searchKey := key + `="`
	start := strings.Index(header, searchKey)
	if start == -1 {
		return ""
	}
	start += len(searchKey)

	end := strings.Index(header[start:], `"`)
	if end == -1 {
		return ""
	}

	return header[start : start+end]
}

func parseWwwAuthenticateManual(header string) (realm, service, scope string) {
	realm = extractValue(header, "realm")
	service = extractValue(header, "service")
	scope = extractValue(header, "scope")
	return
}
