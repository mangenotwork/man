package test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCase1(t *testing.T) {
	name := "Bob"
	age := 10

	assert.Equal(t, "Bob", name)
	assert.Equal(t, 20, age, "年龄不相等")
}

// 简单使用
func TestSimpleRand(t *testing.T) {
	t.Log("start ...")
	assert := assert.New(t)
	assert.Equal(1, 1)
	assert.NotEqual(1, 2)
	assert.NotNil("123")
	assert.IsType([]string{}, []string{""})

	assert.Contains("Hello World", "World")
	assert.Contains(map[string]string{"Hello": "World"}, "Hello")
	assert.Contains([]string{"Hello", "World"}, "Hello")
	assert.True(true)
	assert.True(false)
	t.Log("next ...")
	var s []string
	assert.Empty(s)
	assert.Nil(s)
	t.Log("end ...")
}
