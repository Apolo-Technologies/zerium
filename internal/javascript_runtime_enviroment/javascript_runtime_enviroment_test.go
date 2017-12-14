// Copyright 2017 The zerium Authors
// This file is part of the zerium library.
//
// The zerium library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The zerium library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the zerium library. If not, see <http://www.gnu.org/licenses/>.

package javascript_runtime_enviroment

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/robertkrimen/otto"
)

type testNativeObjectBinding struct{}

type msg struct {
	Msg string
}

func (no *testNativeObjectBinding) TestMethod(call otto.FunctionCall) otto.Value {
	m, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	v, _ := call.Otto.ToValue(&msg{m})
	return v
}

func newWithTestJS(t *testing.T, testjs string) (*JAVASCRIPT_RUNTIME_ENVIROMENT, string) {
	dir, err := ioutil.TempDir("", "javascript_runtime_enviroment-test")
	if err != nil {
		t.Fatal("cannot create temporary directory:", err)
	}
	if testjs != "" {
		if err := ioutil.WriteFile(path.Join(dir, "test.js"), []byte(testjs), os.ModePerm); err != nil {
			t.Fatal("cannot create test.js:", err)
		}
	}
	return New(dir, os.Stdout), dir
}

func TestExec(t *testing.T) {
	javascript_runtime_enviroment, dir := newWithTestJS(t, `msg = "testMsg"`)
	defer os.RemoveAll(dir)

	err := javascript_runtime_enviroment.Exec("test.js")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	val, err := javascript_runtime_enviroment.Run("msg")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !val.IsString() {
		t.Errorf("expected string value, got %v", val)
	}
	exp := "testMsg"
	got, _ := val.ToString()
	if exp != got {
		t.Errorf("expected '%v', got '%v'", exp, got)
	}
	javascript_runtime_enviroment.Stop(false)
}

func TestNatto(t *testing.T) {
	javascript_runtime_enviroment, dir := newWithTestJS(t, `setTimeout(function(){msg = "testMsg"}, 1);`)
	defer os.RemoveAll(dir)

	err := javascript_runtime_enviroment.Exec("test.js")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	time.Sleep(100 * time.Millisecond)
	val, err := javascript_runtime_enviroment.Run("msg")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !val.IsString() {
		t.Errorf("expected string value, got %v", val)
	}
	exp := "testMsg"
	got, _ := val.ToString()
	if exp != got {
		t.Errorf("expected '%v', got '%v'", exp, got)
	}
	javascript_runtime_enviroment.Stop(false)
}

func TestBind(t *testing.T) {
	javascript_runtime_enviroment := New("", os.Stdout)
	defer javascript_runtime_enviroment.Stop(false)

	javascript_runtime_enviroment.Bind("no", &testNativeObjectBinding{})

	_, err := javascript_runtime_enviroment.Run(`no.TestMethod("testMsg")`)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestLoadScript(t *testing.T) {
	javascript_runtime_enviroment, dir := newWithTestJS(t, `msg = "testMsg"`)
	defer os.RemoveAll(dir)

	_, err := javascript_runtime_enviroment.Run(`loadScript("test.js")`)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	val, err := javascript_runtime_enviroment.Run("msg")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !val.IsString() {
		t.Errorf("expected string value, got %v", val)
	}
	exp := "testMsg"
	got, _ := val.ToString()
	if exp != got {
		t.Errorf("expected '%v', got '%v'", exp, got)
	}
	javascript_runtime_enviroment.Stop(false)
}
