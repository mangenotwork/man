package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"sort"
)

// ColorNode 表示K-D树中的一个节点
type ColorNode struct {
	color RGB
	left  *ColorNode
	right *ColorNode
	axis  int
}

// RGB 表示RGB颜色
type RGB struct {
	R, G, B uint8
}

// ColorCluster 表示一个颜色聚类
type ColorCluster struct {
	Color   RGB
	Count   int
	Percent float64
}

// 计算两个RGB颜色之间的欧氏距离
func colorDistance(c1, c2 RGB) float64 {
	dr := float64(c1.R) - float64(c2.R)
	dg := float64(c1.G) - float64(c2.G)
	db := float64(c1.B) - float64(c2.B)
	return math.Sqrt(dr*dr + dg*dg + db*db)
}

// 将颜色转换为RGB结构
func toRGB(c color.Color) RGB {
	r, g, b, _ := c.RGBA()
	return RGB{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
	}
}

// 创建K-D树
func buildKdtree(points []RGB, depth int) *ColorNode {
	if len(points) == 0 {
		return nil
	}

	axis := depth % 3
	sort.Slice(points, func(i, j int) bool {
		switch axis {
		case 0:
			return points[i].R < points[j].R
		case 1:
			return points[i].G < points[j].G
		default:
			return points[i].B < points[j].B
		}
	})

	mid := len(points) / 2
	node := &ColorNode{
		color: points[mid],
		axis:  axis,
	}

	if len(points) > 1 {
		node.left = buildKdtree(points[:mid], depth+1)
		node.right = buildKdtree(points[mid+1:], depth+1)
	}

	return node
}

// 在K-D树中查找最近邻
func nearestNeighbor(root *ColorNode, target RGB, depth int) RGB {
	if root == nil {
		return RGB{}
	}

	axis := depth % 3
	var nextBranch, oppositeBranch *ColorNode

	if (axis == 0 && target.R < root.color.R) ||
		(axis == 1 && target.G < root.color.G) ||
		(axis == 2 && target.B < root.color.B) {
		nextBranch = root.left
		oppositeBranch = root.right
	} else {
		nextBranch = root.right
		oppositeBranch = root.left
	}

	bestColor := root.color
	bestDist := colorDistance(root.color, target)

	if nextBranch != nil {
		nearest := nearestNeighbor(nextBranch, target, depth+1)
		dist := colorDistance(nearest, target)
		if dist < bestDist {
			bestColor = nearest
			bestDist = dist
		}
	}

	if oppositeBranch != nil {
		nearest := nearestNeighbor(oppositeBranch, target, depth+1)
		dist := colorDistance(nearest, target)
		if dist < bestDist {
			bestColor = nearest
			bestDist = dist
		}
	}

	return bestColor
}

// K-means聚类
func kmeans(points []RGB, k int, maxIterations int) []ColorCluster {
	if len(points) == 0 || k <= 0 {
		return nil
	}

	// 初始化中心点
	centroids := make([]RGB, k)
	for i := 0; i < k && i < len(points); i++ {
		centroids[i] = points[i]
	}

	for iter := 0; iter < maxIterations; iter++ {
		// 分配点到最近的中心
		clusters := make([][]RGB, k)
		for _, p := range points {
			minDist := math.MaxFloat64
			closest := 0
			for i, c := range centroids {
				dist := colorDistance(p, c)
				if dist < minDist {
					minDist = dist
					closest = i
				}
			}
			clusters[closest] = append(clusters[closest], p)
		}

		// 更新中心点
		newCentroids := make([]RGB, k)
		changed := false
		for i, cluster := range clusters {
			if len(cluster) == 0 {
				// 如果某个簇为空，重新选择一个随机点作为中心
				newCentroids[i] = points[i%len(points)]
				changed = true
				continue
			}

			// 计算新的中心
			var sumR, sumG, sumB int
			for _, p := range cluster {
				sumR += int(p.R)
				sumG += int(p.G)
				sumB += int(p.B)
			}
			newR := uint8(sumR / len(cluster))
			newG := uint8(sumG / len(cluster))
			newB := uint8(sumB / len(cluster))

			if newR != centroids[i].R || newG != centroids[i].G || newB != centroids[i].B {
				changed = true
			}
			newCentroids[i] = RGB{R: newR, G: newG, B: newB}
		}

		if !changed {
			break
		}

		centroids = newCentroids
	}

	// 构建K-D树以加速最近邻搜索
	root := buildKdtree(centroids, 0)

	// 统计每个聚类的点数量
	counts := make([]int, k)
	for _, p := range points {
		// 使用K-D树找到最近的聚类中心
		nearest := nearestNeighbor(root, p, 0)

		// 找到最近中心在centroids中的索引
		for i, c := range centroids {
			if c == nearest {
				counts[i]++
				break
			}
		}
	}

	// 生成最终结果
	totalPoints := float64(len(points))
	result := make([]ColorCluster, 0, k)
	for i, c := range centroids {
		if counts[i] > 0 {
			result = append(result, ColorCluster{
				Color:   c,
				Count:   counts[i],
				Percent: float64(counts[i]) / totalPoints * 100,
			})
		}
	}

	// 按数量排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})

	return result
}

// 从图片提取颜色
func extractColors(img image.Image, count int) []ColorCluster {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// 收集所有像素颜色
	pixels := make([]RGB, 0, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			rgb := toRGB(c)
			// 过滤透明或接近白色/黑色的颜色
			if rgb.R < 250 || rgb.G < 250 || rgb.B < 250 {
				if rgb.R > 10 || rgb.G > 10 || rgb.B > 10 {
					pixels = append(pixels, rgb)
				}
			}
		}
	}

	// 使用K-means聚类提取主要颜色
	return kmeans(pixels, count, 30)
}

// 创建颜色样本图片
func createColorSample(clusters []ColorCluster, width, height int) image.Image {
	if len(clusters) == 0 {
		return nil
	}

	sampleWidth := width / len(clusters)
	result := image.NewRGBA(image.Rect(0, 0, width, height))

	for i, cluster := range clusters {
		c := color.RGBA{
			R: cluster.Color.R,
			G: cluster.Color.G,
			B: cluster.Color.B,
			A: 255,
		}
		startX := i * sampleWidth
		endX := startX + sampleWidth
		if i == len(clusters)-1 {
			endX = width // 确保最后一个色块填满剩余空间
		}

		for y := 0; y < height; y++ {
			for x := startX; x < endX; x++ {
				result.Set(x, y, c)
			}
		}
	}

	return result
}

func main() {
	inputPath := "./mv3.jpg"
	outputPath := "./out.png"
	colorCount := 10
	flag.Parse()

	// 打开输入图片
	file, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("无法打开图片: %v", err)
	}
	defer file.Close()

	// 解码图片
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("无法解码图片: %v", err)
	}

	// 确保图片是RGBA格式
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// 提取主要颜色
	clusters := extractColors(rgba, colorCount)

	// 打印结果
	fmt.Println("提取的主要颜色:")
	for i, c := range clusters {
		fmt.Printf("%d. RGB(%d, %d, %d) - 占比: %.2f%%\n", i+1, c.Color.R, c.Color.G, c.Color.B, c.Percent)
	}

	// 创建颜色样本图片
	sampleWidth := 400
	sampleHeight := 100
	sampleImg := createColorSample(clusters, sampleWidth, sampleHeight)

	// 保存结果
	outFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("无法创建输出文件: %v", err)
	}
	defer outFile.Close()

	// 根据文件扩展名选择编码格式
	switch ext := getFileExtension(outputPath); ext {
	case "jpg", "jpeg":
		err = jpeg.Encode(outFile, sampleImg, &jpeg.Options{Quality: 95})
	case "png":
		err = png.Encode(outFile, sampleImg)
	default:
		log.Fatalf("不支持的输出格式: %s", ext)
	}

	if err != nil {
		log.Fatalf("无法保存图片: %v", err)
	}

	fmt.Printf("颜色样本已保存至: %s\n", outputPath)
}

// 获取文件扩展名
func getFileExtension(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i+1:]
		}
	}
	return ""
}
