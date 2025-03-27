package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/gofrs/uuid"
	"io"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	case1()

	case3()

	case4()
}

/*
linux：代表 Linux 操作系统。
darwin：代表 macOS 操作系统。
windows：代表 Windows 操作系统。
freebsd：代表 FreeBSD 操作系统。
openbsd：代表 OpenBSD 操作系统。
*/
func case1() {
	os := runtime.GOOS
	switch os {
	case "linux":
		fmt.Println("当前程序运行在 Linux 操作系统上。")
	case "darwin":
		fmt.Println("当前程序运行在 macOS 操作系统上。")
	case "windows":
		fmt.Println("当前程序运行在 Windows 操作系统上。")
	default:
		fmt.Printf("当前程序运行在未知操作系统上: %s\n", os)
	}
}

// 当前多条件抽象成函数
/*

func (result *transmissionResult) IsFailure() bool {
	return result.statusCode != successResponse && result.statusCode != partialSuccessResponse
}

func (result *transmissionResult) CanRetry() bool {
	if result.IsSuccess() {
		return false
	}

	return result.statusCode == partialSuccessResponse ||
		result.retryAfter != nil ||
		(result.statusCode == requestTimeoutResponse ||
			result.statusCode == serviceUnavailableResponse ||
			result.statusCode == errorResponse ||
			result.statusCode == tooManyRequestsResponse ||
			result.statusCode == tooManyRequestsOverExtendedTimeResponse)
}

func (result *transmissionResult) IsPartialSuccess() bool {
	return result.statusCode == partialSuccessResponse &&
		result.response != nil &&
		result.response.ItemsReceived != result.response.ItemsAccepted
}

func (result *transmissionResult) IsThrottled() bool {
	return result.statusCode == tooManyRequestsResponse ||
		result.statusCode == tooManyRequestsOverExtendedTimeResponse ||
		result.retryAfter != nil
}

*/

// case3 给数组类型添加排序函数
type itemTransmissionResults []int

func (results itemTransmissionResults) Len() int {
	return len(results)
}

func (results itemTransmissionResults) Less(i, j int) bool {
	return results[i] < results[j]
}

func (results itemTransmissionResults) Swap(i, j int) {
	tmp := results[i]
	results[i] = results[j]
	results[j] = tmp
}

func case3() {
	var data itemTransmissionResults = itemTransmissionResults{2, 5, 87, 6, 3, 14, 8}
	fmt.Println(data.Len())
	fmt.Println(data.Less(1, 2))
	fmt.Println(data)
	data.Swap(1, 2)
	fmt.Println(data)
}

// case4 gofrs/uuid 库的使用
type uuidGenerator struct {
	sync.Mutex
	fallbackRand *rand.Rand
	reader       io.Reader
}

var uuidgen *uuidGenerator = newUuidGenerator(crand.Reader)

// newUuidGenerator creates a new uuiGenerator with the specified crypto random reader.
func newUuidGenerator(reader io.Reader) *uuidGenerator {
	// Setup seed for fallback random generator
	var seed int64
	b := make([]byte, 8)
	if _, err := io.ReadFull(reader, b); err == nil {
		seed = int64(binary.BigEndian.Uint64(b))
	} else {
		// Otherwise just use the timestamp
		seed = time.Now().UTC().UnixNano()
	}

	return &uuidGenerator{
		reader:       reader,
		fallbackRand: rand.New(rand.NewSource(seed)),
	}
}

// newUUID generates a new V4 UUID
func (gen *uuidGenerator) newUUID() uuid.UUID {
	//call the standard generator
	u, err := uuid.NewV4()
	//err will be either EOF or unexpected EOF
	if err != nil {
		gen.fallback(&u)
	}

	return u
}

// fallback populates the specified UUID with the standard library's PRNG
func (gen *uuidGenerator) fallback(u *uuid.UUID) {
	gen.Lock()
	defer gen.Unlock()
	// This does not fail as per documentation
	gen.fallbackRand.Read(u[:])
	u.SetVersion(uuid.V4)
	u.SetVariant(uuid.VariantRFC4122)
}

// newUUID generates a new V4 UUID
func newUUID() uuid.UUID {
	return uuidgen.newUUID()
}

func case4() {
	var start sync.WaitGroup
	var finish sync.WaitGroup

	start.Add(1)

	goroutines := 250
	uuidsPerRoutine := 10
	results := make(chan string, 100)

	// Start normal set of UUID generation:
	for i := 0; i < goroutines; i++ {
		finish.Add(1)
		go func() {
			defer finish.Done()
			start.Wait()
			for t := 0; t < uuidsPerRoutine; t++ {
				results <- newUUID().String()
			}
		}()
	}

	// Start broken set of UUID generation
	brokenGen := newUuidGenerator(&brokenReader{})
	for i := 0; i < goroutines; i++ {
		finish.Add(1)
		go func() {
			defer finish.Done()
			start.Wait()
			for t := 0; t < uuidsPerRoutine; t++ {
				results <- brokenGen.newUUID().String()
			}
		}()
	}

	// Close the channel when all the goroutines have exited
	go func() {
		finish.Wait()
		close(results)
	}()

	used := make(map[string]bool)
	start.Done()
	for id := range results {
		fmt.Println(id)
		if _, ok := used[id]; ok {
			fmt.Printf("UUID was generated twice: %s", id)
		}

		used[id] = true
	}

}

type brokenReader struct{}

func (reader *brokenReader) Read(b []byte) (int, error) {
	return 0, io.EOF
}
