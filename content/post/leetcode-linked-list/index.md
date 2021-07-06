---
title: "Leetcode Linked List"
description: 
date: 2021-07-06T08:18:24+08:00
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

## [206. 反转链表](https://leetcode-cn.com/problems/reverse-linked-list/)

给你单链表的头节点 head ，请你反转链表，并返回反转后的链表。

示例 1：

![img](https://assets.leetcode.com/uploads/2021/02/19/rev1ex1.jpg)

输入：head = [1,2,3,4,5]
输出：[5,4,3,2,1]

```go
/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
```



### 题解

**递归**

```go
func reverseList(head *ListNode) *ListNode {
    if head == nil {
        return nil
    }
    if head.Next == nil {
        return head
    }
    last := reverseList(head.Next)
    head.Next.Next = head
    head.Next = nil
    return last
}

```



**迭代**

```go
func reverseList(head *ListNode) *ListNode {
    var pre *ListNode
    cur, next := head, head
    for cur != nil {
        next = cur.Next
        cur.Next = pre
        pre = cur
        cur = next
    }
    return pre
}
```



## [92. 反转链表 II](https://leetcode-cn.com/problems/reverse-linked-list-ii/)

给你单链表的头指针 head 和两个整数 left 和 right ，其中 left <= right 。请你反转从位置 left 到位置 right 的链表节点，返回 反转后的链表 。

示例 1：

![img](https://assets.leetcode.com/uploads/2021/02/19/rev2ex2.jpg)


输入：head = [1,2,3,4,5], left = 2, right = 4
输出：[1,4,3,2,5]

### 题解

```go
func reverseBetween(head *ListNode, left int, right int) *ListNode {
    if left == 1 {
        return reverseN(head, right)
    }
    head.Next = reverseBetween(head.Next, left-1, right-1)
    return head
}

var next *ListNode

// 反转链表的前N个节点
func reverseN(head *ListNode, right int) *ListNode {
    if right == 1 {
        next = head.Next
        return head
    }

    last := reverseN(head.Next, right-1)
    head.Next.Next = head
    head.Next = next

    return last

}
```

