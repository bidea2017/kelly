# 路由
1. 路由组
1. middleware（兼容标准库Context）
1. classic模式

## HandlerFunc
兼容http.Handler

## 错误处理
1. 404
1. 405

# 回调Context
### Context数据

1. 基于标准库Context和map
2. 兼容标准库Context(**SetRequest**)

### render
1. json
1. xml
1. text
1. raw data
1. html/template

---

1. 设置header
1. 设置cookie
1. 重定向

### request
1. 读取path变量
1. 读取query变量
1. 读取form变量
1. 读取cookie
1. 读取header

# 内建middleware
