package database

import (
	"github.com/ronmount/ozon_go/internal/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetMD5Hash(t *testing.T) {
	t.Log("When \"123\"")
	expected := "202cb962ac59075b964b07152d234b70"
	actual := GetMD5Hash("123")
	require.Equal(t, expected, actual, "Should be equal")

	t.Log("When \"hello world\"")
	expected = "5eb63bbbe01eeed093cb22bb8f5acdc3"
	actual = GetMD5Hash("hello world")
	require.Equal(t, expected, actual, "Should be equal")
}

func TestMemoryStorage_AddURL(t *testing.T) {
	ms, _ := NewMemoryStorage()

	t.Log("When added first arg")
	link := "https://ozon.ru/"
	res, err := ms.AddURL(link)
	hash := GetMD5Hash(res.(models.Link).ShortLink)
	require.NoError(t, err, "Should be no error")
	require.True(t, len(ms.shortAsKey) == 1, "Should be true")
	require.Equal(t, res.(models.Link).FullLink, link, "Should be equal")
	require.Equal(t, res.(models.Link).FullLink, ms.shortAsKey[hash], "Should be equal")

	t.Log("When added second arg")
	link = "https://vk.com/ronmount"
	res, err = ms.AddURL(link)
	hash = GetMD5Hash(res.(models.Link).ShortLink)
	require.NoError(t, err, "Should be no error")
	require.True(t, len(ms.shortAsKey) == 2, "Should be true")
	require.Equal(t, res.(models.Link).FullLink, link, "Should be equal")
	require.Equal(t, res.(models.Link).FullLink, ms.shortAsKey[hash], "Should be equal")
}

func TestMemoryStorage_GetURL(t *testing.T) {
	ms, _ := NewMemoryStorage()
	links := []string{"https://ozon.ru/", "https://vk.com/ronmount"}
	data := make(map[string]string)
	for _, link := range links {
		res, _ := ms.AddURL(link)
		data[res.(models.Link).FullLink] = res.(models.Link).ShortLink
	}

	t.Logf("When %d links in storage", len(links))
	for k, v := range data {
		res, err := ms.GetURL(v)
		require.NoError(t, err, "Should be no error")
		require.Equal(t, k, res.(models.Link).FullLink, "Should be equal")
	}

	t.Log("When link not in storage")
	_, err := ms.GetURL("wrong_link")
	require.Error(t, err, "Should be error")
}
