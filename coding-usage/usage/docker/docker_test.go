package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/moby/moby/client"
)

var (
	cli *client.Client
	ctx = context.Background()
)

func beforeAfter(caseName string) func() {
	if cli == nil {
		cli, _ = client.New(client.FromEnv, client.WithAPIVersionFromEnv())
	}

	return func() {
		defer func() {
			if cli != nil {
				cli.Close()
			}
		}()
		fmt.Println(caseName, "用例结束...")
	}
}

func Test_docker_image_00(t *testing.T) {
	defer beforeAfter("Test_docker_image_00")()

	images, _ := cli.ImageList(ctx, client.ImageListOptions{All: true, Manifests: true})
	for _, is := range images.Items {
		iir, _ := cli.ImageInspect(ctx, is.ID)
		fmt.Println(time.Unix(is.Created, 0).Format("2006-01-02 15:04:05"), is.Size/1024/1024, iir.Os+"/"+iir.Architecture, is.RepoTags, is.RepoDigests)
		// ghcr.nju.edu.cn/szemeng76/lunatv:latest
		// fullName := is.RepoTags[0]
	}
}

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

// 容器
// result, err := cli.ContainerList(context.Background(), client.ContainerListOptions{
// 	All: true,
// })
// if err != nil {
// 	panic(err)
// }

// // Print each container's ID, status and the image it was created from.
// fmt.Printf("%s  %-22s  %s\n", "ID", "STATUS", "IMAGE")
// for _, ctr := range result.Items {
// 	fmt.Printf("%s  %-22s  %s\n", ctr.ID, ctr.Status, ctr.Image)
// }

// }
