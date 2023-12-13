package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jaswdr/faker"
)

type Node struct {
	ID    int
	Name  string
	Form  string // "circle", "rect", "square", "ellipse", "round-rect", "rhombus"
	Links []*Node
}

func main() {

	BinaryStart()
	graph()
	setContent()
	Worker()

	rp := NewReverseProxy("hugo", "1313")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(rp.ReverseProxy)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from API"))
	})

	r.Handle("/*", nextHandler)
	http.ListenAndServe(":8080", r)

}

func BinaryStart() {

	cicle := 50

	tr := GenerateTree(cicle)

	time.Sleep(2 * time.Second)

	for i := 0; i < cicle; i++ {

		tr.PrintThree()

		fmt.Println(ResultString)
		printDocument(ResultString)
		ResultString = ""

		time.Sleep(2 * time.Second)
	}
}

func printDocument(str string) {

	SourcePath := "/home/ars/go/hugopro/hugo/content/tasks/blank/binary.md"
	//SourcePath := "/src/static/tasks/blank/binary.md"

	data, err := os.ReadFile(SourcePath)
	if err != nil {
		log.Fatal(err)
	}

	redData := fmt.Sprintf(string(data), str)

	DistPath := "/home/ars/go/hugopro/hugo/content/tasks/binary.md"
	//DistPath := "/src/static/tasks/graph.md"
	err = os.WriteFile(DistPath, []byte(redData), 0777)

	if err != nil {
		log.Println(err)
	}
}

var graphString string

func graph() {

	var currentNode *Node
	var rootNode *Node

	forms := []string{"circle", "rect", "square", "ellipse", "round-rect", "rhombus"}
	rand.Seed(time.Now().UnixNano())

	fake := faker.New()
	p := fake.Person()

	cnt := 0

	for b := 0; b < 20; b++ {

		for i := 0; i < rand.Intn(20)+5; i++ {

			graphString = ""

			newNode := newNode(i, p.FirstName(), forms[rand.Intn(len(forms)-1)])

			if cnt > 3 {
				cnt = 0
			}

			if cnt == 0 && i == 0 {

				currentNode = newNode
				rootNode = newNode
				continue
			}

			currentNode.Links = append(currentNode.Links, newNode)

			if rand.Intn(50) < 25 {

				currentNode = newNode
			}
		}

		cnt++

		time.Sleep(2 * time.Second)

		parseNode(*rootNode, 0, Node{})
		fmt.Println(graphString)

		//SourcePath := "/home/ars/go/hugopro/hugo/content/tasks/g.md"
		SourcePath := "/app/static/tasks/g.md"

		data, err := os.ReadFile(SourcePath)
		if err != nil {
			log.Fatal(err)
		}

		redData := fmt.Sprintf(string(data), graphString)

		//DistPath := "/home/ars/go/hugopro/hugo/content/tasks/graph.md"
		DistPath := "/app/static/tasks/graph.md"

		err = os.WriteFile(DistPath, []byte(redData), 0777)

		if err != nil {
			log.Println(err)
		}

	}
}

func newNode(id int, name, form string) *Node {
	return &Node{
		ID:    id,
		Name:  name,
		Form:  form,
		Links: []*Node{},
	}
}

func parseNode(node Node, count int, parent Node) {

	if parent.ID != 0 {

		st := fmt.Sprintf("%s --> %s \n", parent.Name, node.Name)

		graphString += st
	}

	count++

	for _, link := range node.Links {

		parseNode(*link, count, node)
	}
}

var content string

type ReverseProxy struct {
	host string
	port string
}

func NewReverseProxy(host, port string) *ReverseProxy {
	return &ReverseProxy{
		host: host,
		port: port,
	}
}

func (rp *ReverseProxy) ReverseProxy(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/api" || strings.HasPrefix(r.URL.Path, "/api/") {

			next.ServeHTTP(w, r)

		} else {

			u, _ := url.Parse("http://hugo:1313")
			r.URL.Scheme = u.Scheme
			r.URL.Host = u.Host
			r.URL.Path = u.Path
			proxy := httputil.NewSingleHostReverseProxy(u)
			proxy.ServeHTTP(w, r)
		}
	})
}

func setContent() {

	content = `
---
menu:
    before:
        name: tasks
        weight: 5
title: Обновление данных в реальном времени
---

# Задача: Обновление данных в реальном времени

Напишите воркер, который будет обновлять данные в реальном времени, на текущей странице.
Текст данной задачи менять нельзя, только время и счетчик.

Файл данной страницы: /app/static/tasks/_index.md	

Должен меняться счетчик и время:

Текущее время: %s

Счетчик: %d

## Критерии приемки:
- [ ] Воркер должен обновлять данные каждые 5 секунд
- [ ] Счетчик должен увеличиваться на 1 каждые 5 секунд
- [ ] Время должно обновляться каждые 5 секунд`

}

func WorkerTest() {
	t := time.NewTicker(5 * time.Second)
	var b byte = 0
	for {
		select {
		case <-t.C:
			err := os.WriteFile("/app/static/_index.md", []byte(fmt.Sprintf(content, b)), 0644)
			if err != nil {
				log.Println(err)
			}
			b++
		}
	}
}

func Worker() {

	t := time.NewTicker(5 * time.Second)
	var counter int
	for {
		select {
		case <-t.C:

			currentTime := time.Now().Format("2006-01-02 15:04:05")

			path := "/app/static/tasks/_index.md"
			//path := "../hugo/content/tasks/_index.md"

			err := os.WriteFile(path, []byte(fmt.Sprintf(content, currentTime, counter)), 0777)

			if err != nil {
				log.Println(err)
			}

			counter++
		}
	}
}

// =========================

type TreeNode struct {
	Key    int
	Height int
	Left   *TreeNode
	Right  *TreeNode
}

type AVLTree struct {
	Root *TreeNode
}

var ResultString string

func (t *AVLTree) PrintThree() {

	ParseLeaf(t.Root, 0)
}

func ParseLeaf(n *TreeNode, parent int) {

	if parent != 0 {
		v := fmt.Sprintf("%d --> %d \n", parent, n.Key)
		ResultString += v
	}

	if n.Left != nil {
		ParseLeaf(n.Left, n.Key)
	}

	if n.Right != nil {
		ParseLeaf(n.Right, n.Key)
	}
}

func NewNode(key int) *TreeNode {
	return &TreeNode{Key: key, Height: 1}
}

func (t *AVLTree) Insert(key int) {
	t.Root = insert(t.Root, key)
}

func (t *AVLTree) ToMermaid() string {
	return ""
}

func height(node *TreeNode) int {
	if node == nil {
		return 0
	}
	return node.Height
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func updateHeight(node *TreeNode) {
	node.Height = 1 + max(height(node.Left), height(node.Right))
}

func getBalance(node *TreeNode) int {
	if node == nil {
		return 0
	}
	return height(node.Left) - height(node.Right)
}

func leftRotate(x *TreeNode) *TreeNode {
	y := x.Right
	x.Right = y.Left
	y.Left = x
	updateHeight(x)
	updateHeight(y)
	return y
}

func rightRotate(y *TreeNode) *TreeNode {
	x := y.Left
	y.Left = x.Right
	x.Right = y
	updateHeight(y)
	updateHeight(x)
	return x
}

func insert(node *TreeNode, key int) *TreeNode {
	if node == nil {
		return NewNode(key)
	}
	if key < node.Key {
		node.Left = insert(node.Left, key)
	} else if key > node.Key {
		node.Right = insert(node.Right, key)
	} else {
		return node
	}
	updateHeight(node)
	balance := getBalance(node)
	if balance > 1 && key < node.Left.Key {
		return rightRotate(node)
	}
	if balance < -1 && key > node.Right.Key {
		return leftRotate(node)
	}
	if balance > 1 && key > node.Left.Key {
		node.Left = leftRotate(node.Left)
		return rightRotate(node)
	}
	if balance < -1 && key < node.Right.Key {
		node.Right = rightRotate(node.Right)
		return leftRotate(node)
	}
	return node
}

func GenerateTree(count int) *AVLTree {

	tree := &AVLTree{}

	go func() {

		cnt := 0
		for i := 1; i < count; i++ {

			tree.Insert(i)

			if i > 2 {

				time.Sleep(2 * time.Second)
			}

			cnt++

			if cnt > 10 {

				tree.Root.Left = nil
				tree.Root.Right = nil
				//tree = &AVLTree{}
				cnt = 0
			}
		}

	}()

	return tree
}
