---
title: "Binary Tree"
description: 
date: 2021-07-11T04:50:51+08:00
image: 
math: 
license: 
hidden: false
draft: true
categories:
    - leetcode
tags:
    - linked list
    - recursion
---

## [226. 翻转二叉树](https://leetcode-cn.com/problems/invert-binary-tree/)

翻转一棵二叉树。

示例：

输入：

```shell
     4
   /   \
  2     7
 / \   / \
1   3 6   9
```

输出：

```shell
     4
   /   \
  7     2
 / \   / \
9   6 3   1
```
### 题解

```go
func invertTree(root *TreeNode) *TreeNode {
    if root == nil {
        return root
    }
    tmp := root.Left
    root.Left = root.Right
    root.Right = tmp
    invertTree(root.Left)
    invertTree(root.Right)
    
    // or 后序
    // invertTree(root.Left)
    // invertTree(root.Right)
    // mp := root.Left
    // root.Left = root.Right
    // root.Right = tmp
    
    return root
}
```



## [114. 二叉树展开为链表](https://leetcode-cn.com/problems/flatten-binary-tree-to-linked-list/)

给你二叉树的根结点 root ，请你将它展开为一个单链表：

展开后的单链表应该同样使用 TreeNode ，其中 right 子指针指向链表中下一个结点，而左子指针始终为 null 。
展开后的单链表应该与二叉树 先序遍历 顺序相同。


![img](https://assets.leetcode.com/uploads/2021/01/14/flaten.jpg)



### 题解

```go
func flatten(root *TreeNode)  {
    if root == nil {
        return 
    }
   
    flatten(root.Left)
    flatten(root.Right)
    right := root.Right
    left := root.Left
    root.Left = nil
    root.Right = left
    p := root
    for p.Right != nil {
        p = p.Right
    } 
    p.Right = right
}
```



## [116. 填充每个节点的下一个右侧节点指针](https://leetcode-cn.com/problems/populating-next-right-pointers-in-each-node/)

给定一个 完美二叉树 ，其所有叶子节点都在同一层，每个父节点都有两个子节点。二叉树定义如下：

```go
struct Node {
  int val;
  Node *left;
  Node *right;
  Node *next;
}
```


填充它的每个 next 指针，让这个指针指向其下一个右侧节点。如果找不到下一个右侧节点，则将 next 指针设置为 NULL。

初始状态下，所有 next 指针都被设置为 NULL。

**进阶：**

- 你只能使用常量级额外空间。
- 使用递归解题也符合要求，本题中递归程序占用的栈空间不算做额外的空间复杂度。

**示例：**

![img](https://assets.leetcode.com/uploads/2019/02/14/116_sample.png)

```shell
输入：root = [1,2,3,4,5,6,7]
输出：[1,#,2,3,#,4,5,6,7,#]
解释：给定二叉树如图 A 所示，你的函数应该填充它的每个 next 指针，以指向其下一个右侧节点，如图 B 所示。序列化的输出按层序遍历排列，同一层节点由 next 指针连接，'#' 标志着每一层的结束。
```

### 题解

```go
func connect(root *Node) *Node {
	if root == nil {
        return root
    }
    connectTwo(root.Left, root.Right)
    return root
}

func connectTwo(node1, node2 *Node) {
    if node1 == nil{
        return
    }
    node1.Next = node2
    connectTwo(node1.Left, node1.Right)
    connectTwo(node2.Left, node2.Right)
    connectTwo(node1.Right, node2.Left)
}
```

