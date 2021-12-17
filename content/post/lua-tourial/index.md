---
title: "Lua Tourial"
description: 
date: 2021-08-20T10:55:34+08:00
image: 
math: 
license: 
hidden: false
draft: true
category:
    -: lua
tag:
    -: lua
---

## 数据类型

- nil
- boolean
- number
- string
- function
- userdata
- thread
- table

### boolean

包括

- true
- false

lua将false以及nil视为fales，其余均是true（包括数字0）

```lua
if 0 then
    print("ok")
else
    print("not ok")
end

-- output
-- ok
```



### number

lua只有一种数字类型——number（双精度double）

```lua
print(2+0.2)
print(2+3)
print(2+2e+1)
print(2e-1)

-- output
-- 2.2
-- 5
-- 22
-- 0.2
```



### string

单行字符串使用单引号或者多引号表示

多行字符串使用`[[]]`表示

字符串拼接使用`..`连接

使用`+`可以对字符串进行算术操作，如果字符串无法转换成数字将会报错

使用`#`来计算字符串长度

```lua
str = [[
    I
    am
    BatMan
]]
print(str .. "\t!")
print(type("2"+6))
print("2e+1"+2)
print(#str)

-- output
--     I
--     am
--     BatMan
--         !
-- number
-- 22
-- 24
```



### table

table

```lua
tbl = {"apple", "pear", "orange", "grape"}
print(tbl[1])
tbl["foo"] = "bar"
for key, val in pairs(tbl) do
    print(key .. " : " .. val)
end

-- output
-- apple
-- 1 : apple
-- 2 : pear
-- 3 : orange
-- 4 : grape
-- foo : bar
```

