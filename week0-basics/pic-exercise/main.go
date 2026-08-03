package main

import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {
	// slice externo: dy linhas
	imagem := make([][]uint8, dy)

	for y := 0; y < dy; y++ {
		// cada linha é seu próprio slice, com dx colunas
		imagem[y] = make([]uint8, dx)

		for x := 0; x < dx; x++ {
			imagem[y][x] = uint8((x + y) / 2)
		}
	}

	return imagem
}

func main() {
	pic.Show(Pic)
}
