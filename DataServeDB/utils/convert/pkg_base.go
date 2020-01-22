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

import "fmt"

// declarions


//public type decl
type ConversionClass int

const (
	Strict ConversionClass = iota 	//Only if type is same.

	Lossless //Only if conversion doesn't loose data. e.g. int to and from string will not lose data, hence, it
	// will not error under Lossless conversion type.

	Weak	//It will convert in most cases but to and from type may not result is same value.
	// For example float might lose data converting to and from string.
)

// public generic functions

// private generic functions

func capByLimits(f float64, min float64, max float64) float64 {
	if f >= min && f <= max {
		return f
	} else if f >= min {
		return max
	} else {
		return min
	}
}

func losslessError(cc ConversionClass, i interface{}, toTypeName string) error {
	if cc == Lossless {
		return fmt.Errorf("could not convert value '%#v' of type '%T' to type '%v' due to lossless type conversion rule", i, i, toTypeName)
	}
	return nil
}

func strictError(cc ConversionClass, i interface{}, toTypeName string) error {
	if cc == Strict {
		return fmt.Errorf("could not convert value '%#v' of type '%T' to type '%v' due to strict type conversion rule", i, i, toTypeName)
	}
	return nil
}

