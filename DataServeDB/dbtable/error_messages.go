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

package dbtable

import (
	"fmt"
)

//This file helps reuse code for error replies.

func errRplFieldDoesNotExist(fieldName string) error {
	return fmt.Errorf("field '%s' does not exist", fieldName)
}

func errRplFieldNameAlreadyExist(fieldName string) error {
	return fmt.Errorf("field name '%s' already exists", fieldName)
}

func errRplRowDataConversion(fieldName string, conversionError error) error {
	return fmt.Errorf("error occured for field '%v': %v", fieldName, conversionError)
}
