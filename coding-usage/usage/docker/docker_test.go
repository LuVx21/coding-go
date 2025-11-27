package main

// import (
// 	"context"
// 	"fmt"
// 	"github.com/docker/docker/api/types/container"
// 	"github.com/docker/docker/api/types/image"
// 	"github.com/docker/docker/client"
// 	"github.com/luvx21/coding-go/coding-common/common_x"
// 	"strings"
// 	"testing"
// )

// var cli *client.Client
// var ctx = context.Background()

// func beforeAfter(caseName string) func() {
// 	if cli == nil {
// 		cli, _ = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// 	}

// 	return func() {
// 		fmt.Println(caseName, "用例结束...")
// 	}
// }

// func Test_docker_image_00(t *testing.T) {
// 	defer beforeAfter("Test_docker_image_00")()

// 	images, _ := cli.ImageList(ctx, image.ListOptions{All: true})
// 	for _, i := range images {
// 		for _, repository := range i.RepoTags {
// 			if strings.Contains(repository, "luvx") || strings.Contains(repository, "none") {
// 				continue
// 			}
// 			fmt.Println("拉取/更新镜像:", repository)
// 			_, _ = cli.ImagePull(ctx, repository, image.PullOptions{})
// 		}
// 	}
// }
// func Test_image_backup(t *testing.T) {
// 	defer beforeAfter("Test_image_backup")()

// 	imageId := ""
// 	images, _ := cli.ImageList(ctx, image.ListOptions{All: true})
// 	for _, i := range images {
// 		if !strings.Contains(i.ID, imageId) {
// 			continue
// 		}
// 		for _, repository := range i.RepoTags {
// 			split := strings.Split(repository, ":")
// 			_image := split[0]
// 			version := common_x.IfThen(len(split) > 1, split[1], "latest")
// 			fmt.Println(_image, version)

// 			for _, registry := range []string{"registry.cn-shanghai.aliyuncs.com", "ccr.ccs.tencentyun.com"} {
// 				nImage := registry + "/luvx21/" + strings.ReplaceAll(_image, "/", "_")
// 				nn := nImage + ":" + version
// 				fmt.Println("备份为->", nn)
// 				_ = cli.ImageTag(ctx, imageId, nn)
// 				_, _ = cli.ImagePush(ctx, nn, image.PushOptions{})
// 				_, _ = cli.ImageRemove(ctx, nn, image.RemoveOptions{})
// 			}
// 		}
// 	}
// }

// func Test_docker_container_00(t *testing.T) {
// 	defer beforeAfter("Test_docker_container_00")()

// 	containers, _ := cli.ContainerList(ctx, container.ListOptions{All: true})
// 	for _, c := range containers {
// 		fmt.Println(c.Image, c.Ports, c.Names)
// 	}
// }
