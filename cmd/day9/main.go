package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/rawbits2010/AoC24/internal/inputhandler"
)

func main() {

	lines := inputhandler.ReadInput()

	disk, fileCatalog, err := decipherDiskLayout(lines[0])
	if err != nil {
		log.Fatal(err)
	}

	compactedDiskP1 := compactingP1(disk)
	resultP1 := generateChecksum(compactedDiskP1)

	compactedDiskP2 := compactingP2(fileCatalog)
	resultP2 := generateChecksum(compactedDiskP2)

	fmt.Printf("Result - Part 1: %d, Part 2: %d ", resultP1, resultP2)
}

type DiskRange struct {
	DiskPos int
	IsFree  bool
	FileIdx int
	Size    int
}

func decipherDiskLayout(line string) ([]int, []DiskRange, error) {

	layoutValues := make([]int, len(line))

	var diskLength int
	for currIdx := 0; currIdx < len(line); currIdx++ {

		currValStr := line[currIdx : currIdx+1]
		val, err := strconv.Atoi(currValStr)
		if err != nil {
			return nil, nil, fmt.Errorf("error converting disk layout value '%s'", currValStr)
		}

		layoutValues[currIdx] = val
		diskLength += val
	}

	fileCatalog := make([]DiskRange, 0, len(layoutValues))
	disk := make([]int, diskLength)
	var currDiskPos int
	for idx, spaceLength := range layoutValues {

		rangeEntry := DiskRange{
			DiskPos: currDiskPos,
			Size:    spaceLength,
		}

		if idx%2 == 0 {
			fileIdx := idx / 2

			rangeEntry.IsFree = false
			rangeEntry.FileIdx = fileIdx

			for ; spaceLength > 0; spaceLength-- {

				disk[currDiskPos] = fileIdx

				currDiskPos++
			}
		} else {

			rangeEntry.IsFree = true
			rangeEntry.FileIdx = -1

			for ; spaceLength > 0; spaceLength-- {

				disk[currDiskPos] = -1

				currDiskPos++
			}
		}

		fileCatalog = append(fileCatalog, rangeEntry)
	}

	return disk, fileCatalog, nil
}

func compactingP1(disk []int) []int {

	for currDiskIdx := 0; currDiskIdx < len(disk); currDiskIdx++ {

		if disk[currDiskIdx] != -1 {
			continue
		}

		var lastEntry int
		lastEntryIdx := len(disk) - 1
		for {
			lastEntry = disk[lastEntryIdx]

			if lastEntry == -1 {
				lastEntryIdx--
				continue
			}

			break
		}

		if lastEntryIdx <= currDiskIdx {
			disk = disk[:currDiskIdx+1]
			break
		}

		disk[currDiskIdx] = lastEntry
		disk = disk[:lastEntryIdx]
	}

	return disk
}

func compactingP2(fileCatalog []DiskRange) []int {

	for currEntryIdx := len(fileCatalog) - 1; currEntryIdx >= 0; currEntryIdx-- {

		if fileCatalog[currEntryIdx].IsFree {
			continue
		}

		for currFreeSpace := 0; currFreeSpace < currEntryIdx; currFreeSpace++ {

			if fileCatalog[currFreeSpace].IsFree {
				if fileCatalog[currFreeSpace].Size >= fileCatalog[currEntryIdx].Size {

					extraFreeSpace := fileCatalog[currFreeSpace].Size - fileCatalog[currEntryIdx].Size

					fileCatalog[currFreeSpace].FileIdx = fileCatalog[currEntryIdx].FileIdx
					fileCatalog[currFreeSpace].Size = fileCatalog[currEntryIdx].Size
					fileCatalog[currFreeSpace].IsFree = false

					fileCatalog[currEntryIdx].FileIdx = -1
					fileCatalog[currEntryIdx].IsFree = true

					if extraFreeSpace > 0 {

						tempEntry := DiskRange{
							DiskPos: fileCatalog[currFreeSpace].DiskPos + fileCatalog[currFreeSpace].Size,
							IsFree:  true,
							FileIdx: -1,
							Size:    extraFreeSpace,
						}

						tempCatalog := make([]DiskRange, 0, len(fileCatalog)+1)
						tempCatalog = append(tempCatalog, fileCatalog[:currFreeSpace+1]...)
						tempCatalog = append(tempCatalog, tempEntry)
						tempCatalog = append(tempCatalog, fileCatalog[currFreeSpace+1:]...)

						fileCatalog = tempCatalog
					}

					break
				}
			}
		}
	}

	return materializeDisk(fileCatalog)
}

func materializeDisk(fileCatalog []DiskRange) []int {

	var diskSize int
	for _, entry := range fileCatalog {
		diskSize += entry.Size
	}

	disk := make([]int, diskSize)
	for idx := 0; idx < len(disk); idx++ {
		disk[idx] = -1
	}

	for _, entry := range fileCatalog {

		if entry.IsFree {
			continue
		}

		for currDiskPos := entry.DiskPos; currDiskPos < entry.DiskPos+entry.Size; currDiskPos++ {
			disk[currDiskPos] = entry.FileIdx
		}
	}

	return disk
}

func generateChecksum(disk []int) uint64 {

	var checksum uint64
	for blockIdx := 0; blockIdx < len(disk); blockIdx++ {

		if disk[blockIdx] == -1 {
			continue
		}

		checksum += uint64(blockIdx * disk[blockIdx])
	}

	return checksum
}

func printDisk(disk []int) {
	for _, v := range disk {
		if v == -1 {
			fmt.Printf(".")
		} else {
			fmt.Printf("%d", v)
		}
	}
	fmt.Print("\n")
}
