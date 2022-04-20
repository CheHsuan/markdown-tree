package mdtree

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	exampleDirSuffix = "-example"
	maxLevel         = 6
)

type markdownTree struct {
	root  string
	nodes map[string]interface{}
}

func NewMarkdownTree(root string) *markdownTree {
	return &markdownTree{
		root:  root,
		nodes: map[string]interface{}{},
	}
}

func (tree *markdownTree) InsertNode(path string, isDir bool) bool {
	parts := strings.Split(path, "/")

	// skip cases that files in the root dir, should be at least {root_dir}/XXX/YYY
	if len(parts) < 2 {
		return false
	}

	// skip files in example dir
	if strings.Contains(filepath.Dir(path), exampleDirSuffix) {
		return false
	}

	// starting from 2 is because level 1 is root dir
	tree.insertNode(2, parts, isDir, tree.nodes)
	return true
}

func (tree *markdownTree) insertNode(level int, parts []string, isDir bool, nodes map[string]interface{}) {
	if len(parts) == 2 || level >= maxLevel {
		if _, ok := nodes[parts[0]]; !ok {
			nodes[parts[0]] = map[string]interface{}{}
		}
		if !isDir || strings.Contains(parts[1], exampleDirSuffix) || level >= maxLevel {
			nodes[parts[0]].(map[string]interface{})[parts[1]] = getFileNameWithoutExtension(parts[1])
		} else if isDir {
			nodes[parts[0]].(map[string]interface{})[parts[1]] = map[string]interface{}{}
		}
		return
	}

	tree.insertNode(level+1, parts[1:], isDir, nodes[parts[0]].(map[string]interface{}))
}

func (tree *markdownTree) GenerateMarkdown(output string) error {
	fd, err := os.Create(output)
	if err != nil {
		return err
	}
	defer fd.Close()

	w := bufio.NewWriter(fd)
	defer w.Flush()

	tree.generateMarkdown(w, 1, tree.nodes)

	return nil
}

func (tree *markdownTree) generateMarkdown(w *bufio.Writer, level int, nodes map[string]interface{}) {
	if level == 1 {
		w.WriteString(fmt.Sprintf("%s %s\n", getNumberSign(level), tree.root))
		tree.generateMarkdown(w, level+1, nodes)
		return
	}

	for k, f := range nodes {
		if _, ok := f.(string); ok {
			w.WriteString(fmt.Sprintf("- %s\n", k))
		}
	}

	for k, f := range nodes {
		if fmap, ok := f.(map[string]interface{}); ok {
			w.WriteString(fmt.Sprintf("%s %s\n", getNumberSign(level), k))
			tree.generateMarkdown(w, level+1, fmap)
		}
	}
}

func getNumberSign(count int) string {
	str := ""
	for i := 0; i < count; i++ {
		str += "#"
	}
	return str
}

func getFileNameWithoutExtension(fileName string) string {
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}
