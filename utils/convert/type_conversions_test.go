// Copyright (c) 2019 Advanced Computing Labs DMCC

/*
	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
	THE SOFTWARE.
*/

package convert

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestToBoolFunction(t *testing.T) {
	//bool, weak.
	if r, e := ToBool(true, Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != true {
			t.Error("Should be true.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	if r, e := ToBool(false, Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != false {
			t.Error("Should be false.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	//bool, strong.
	if r, e := ToBool(true, Lossless); e == nil {
		fmt.Println("Result: ", r)
		if r != true {
			t.Error("Should be true.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	if r, e := ToBool(false, Lossless); e == nil {
		fmt.Println("Result: ", r)
		if r != false {
			t.Error("Should be false.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	//int, weak
	if r, e := ToBool(int(0), Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != false {
			t.Error("Should be false.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	if r, e := ToBool(int(1), Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != true {
			t.Error("Should be true.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	if r, e := ToBool(int(101), Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != true {
			t.Error("Should be true.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	//int, strong
	if r, e := ToBool(int(0), Lossless); e == nil {
		fmt.Println("Result: ", r)
		t.Error("Should error.")
	} else {
		fmt.Println(e)
	}

	if r, e := ToBool(int(1), Lossless); e == nil {
		fmt.Println("Result: ", r)
		t.Error("Should error.")
	} else {
		fmt.Println(e)
	}

	if r, e := ToBool(int(101), Lossless); e == nil {
		fmt.Println("Result: ", r)
		t.Error("Should error.")
	} else {
		fmt.Println(e)
	}

	//uint, weak
	if r, e := ToBool(uint(0), Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != false {
			t.Error("Should be false.")
		}
	} else {
		t.Error("Should not have error.", e)
	}

	//int64, weak
	if r, e := ToBool(int64(0), Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != false {
			t.Error("Should be false.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	if r, e := ToBool(int64(101), Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != true {
			t.Error("Should be true.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	//int64, strong
	if r, e := ToBool(int64(0), Lossless); e == nil {
		fmt.Println("Result: ", r)
		t.Error("Should error.")
	} else {
		fmt.Println(e)
	}

	if r, e := ToBool(int64(101), Lossless); e == nil {
		fmt.Println("Result: ", r)
		t.Error("Should error.")
	} else {
		fmt.Println(e)
	}

	//float, weak
	if r, e := ToBool(float32(0), Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != false {
			t.Error("Should be false.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	if r, e := ToBool(float32(101), Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != true {
			t.Error("Should be true.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	//int64, strong
	if r, e := ToBool(float32(0), Lossless); e == nil {
		fmt.Println("Result: ", r)
		t.Error("Should error.")
	} else {
		fmt.Println(e)
	}

	if r, e := ToBool(float32(101), Lossless); e == nil {
		fmt.Println("Result: ", r)
		t.Error("Should error.")
	} else {
		fmt.Println(e)
	}

	//string, weak
	if r, e := ToBool("false", Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != false {
			t.Error("Should be false.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	if r, e := ToBool("FALSE", Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != false {
			t.Error("Should be false.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	if r, e := ToBool("true", Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != true {
			t.Error("Should be true.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	if r, e := ToBool("false1", Weak); e == nil {
		fmt.Println("Result: ", r)
		t.Error("Should error.")
	} else {
		fmt.Println(e)
	}

	if r, e := ToBool("true1", Weak); e == nil {
		fmt.Println("Result: ", r)
		t.Error("Should error.")
	} else {
		fmt.Println(e)
	}

	//string, strong
	if r, e := ToBool("false", Lossless); e == nil {
		fmt.Println("Result: ", r)
		t.Error("Should error.")
	} else {
		fmt.Println(e)
	}

	if r, e := ToBool("true", Lossless); e == nil {
		fmt.Println("Result: ", r)
		t.Error("Should error.")
	} else {
		fmt.Println(e)
	}

	//nil, weak
	if r, e := ToBool(nil, Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != false {
			t.Error("Should be false.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	//nil, strong
	if r, e := ToBool(nil, Lossless); e == nil {
		fmt.Println("Result: ", r)
		t.Error("Should error.")
	} else {
		fmt.Println(e)
	}
}

func TestToInt32Function(t *testing.T) {

	//int, strict
	if unsafe.Sizeof(int(0)) == 4 {
	 //TODO:
	}

	if unsafe.Sizeof(int(0)) == 8 {
		if r, e := ToInt32(int(101), Strict); e == nil {
			fmt.Println("Result: ", r)
			t.Error("Should error.")
		} else {
			fmt.Println(e)
		}
	}

	//int, weak
	if r, e := ToInt32(int(0), Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != 0 {
			t.Error("Should be 0.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	if r, e := ToInt32(int(101), Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != 101 {
			t.Error("Should be 0.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	if r, e := ToInt32(int(4294967295), Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != 2147483647 {
			t.Error("Should be 0.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	//int, strong
	if r, e := ToInt32(int(4294967295), Lossless); e == nil {
		fmt.Println("Result: ", r)
		t.Error("Should error.")
	} else {
		fmt.Println(e)
	}

	//int8, weak
	if r, e := ToInt32(int8(0), Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != 0 {
			t.Error("Should be 0.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	if r, e := ToInt32(int8(101), Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != 101 {
			t.Error("Should be 0.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	//unint32, weak.
	if r, e := ToInt32(uint32(4294967295), Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != 2147483647 {
			t.Error("Should be 0.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	//int64, weak.
	if r, e := ToInt32(int64(-9223372036854775808), Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != -2147483648 {
			t.Error("Should be 0.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	if r, e := ToInt32(int64(9223372036854775807), Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != 2147483647 {
			t.Error("Should be 0.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	//float, weak
	if r, e := ToInt32(float32(-3.4e38), Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != -2147483648 {
			t.Error("Should be 0.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	if r, e := ToInt32(float32(3.4e38), Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != 2147483647 {
			t.Error("Should be 0.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	if r, e := ToInt32(float32(-12.5), Weak); e == nil {
		fmt.Println("Result: ", r)
		//if r != 2147483647 {
		//	t.Error("Should be 0.")
		//}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	//string, weak
	if r, e := ToInt32("922337203685477580.7", Weak); e == nil {
		fmt.Println("Result: ", r)
		//if r != 2147483647 {
		//	t.Error("Should be 0.")
		//}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

}

func TestToStringFunction(t *testing.T) {
	//int8, weak
	if r, e := ToString(int8(127), Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != "127" {
			t.Error("Should be 0.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}

	//int64, weak
	if r, e := ToString(int64(9223372036854775807), Weak); e == nil {
		fmt.Println("Result: ", r)
		if r != "9223372036854775807" {
			t.Error("Should be 0.")
		}
	} else {
		fmt.Println(e)
		t.Error("Should not have error.")
	}
}

func TestToIso8601Utc(t *testing.T) {
	dtIso, e := ToIso8601Utc("0900-01-01", Lossless)
	if e == nil {
		fmt.Println(dtIso)
	} else {
		fmt.Println(e)
	}
}
