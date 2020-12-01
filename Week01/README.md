Golang error处理
1.1 Golang 错误处理
Go语言是多返回值的，通常使用bool或者error来确定函数调用成功或者失败
bool：当错误原因只有一种，或者上层应用不关心错误原因，只关心结果时使用
error：错误的原因有很多，且调用者关心错误原因时使用
1.2 error 本质
```
// error 是一个普通的接口
type error interface {
	Error() string
}
```
1.3 error type
```
1.3.1 sentinel error // 哨兵错误
    这一类错误是预定义的，通常使用在一些特定的error处理，如 io包 io.EOF，sql包的 sql.ErrNoRows
    调用者一般需要做等值判断 如 if err == io.EOF { // do somthing}
    这种预定义的错误有时无法避免，但是从设计上似乎有些欠妥
    Go1.3 之后可以使用 if errors.Is(err,io.EOF){ // do something},使用起来更加优雅
    
1.3.2 type error // 自定义错误
    这一类是用户自定义的错误，使用基础的错误可能会遗失一些关键的元数据，所以通过自定义的error来解决此问题
    type MyError struct{
        something string
        msg string
    }
    
    func (e *MyError) Error()string{
        return m.something + e.msg
    }
    
    func (e* MyError)Something() string{
        return e.something
    }
    
    调用者需要使用类型断言 如: if err.(MyError).Something() == "something" {// do something}
    Go1.3 之后可以使用 if errors.As(err,MyError){ // do something},使用起来更加优雅

```
1.4 handle error
```
1.4.1 在不破坏程序完整性，可读性的前提下，应该尽可能的减少对error的处理 如:

// bad case
func biz() error{
    if err:= dao();err!=nil{
        return err
    }
        
    return nilll
} 
    
// good case
func biz() error{
    return dao()
}  
    
1.4.2 error 只应该被处理一次，打日志也算处理
    为了避免error的多次处理，与日志冗余，error只在最顶层处理，下层业务只负责往事抛出
    但是这样出现问题有很难排查，所以需要error带上一些元数据 如：堆栈信息
    引入 github.com/pkg/errors 这个包，多error进行Wrap操作 例子参考作业：work.go

1.4.3 wrap error 的注意事项
    error每Wrap一次都会带上一次堆栈信息，如果每层都wrap的话，会导致error中带入了大量的重复信息
    所以每个error只应该被wrap一次，就是在error第一次产生的地方
    
    error的wrap也只应该发生在业务代码中，不应该在出现在基础库中，原因：队友调用者来说，他是无法甄别基础库是否
    使用了wrap，自己发生错误时是否应该wrap。所以基础库不wrap，可以使用withMessage加入一些元数据但是不包含堆栈信息
    业务代码使用wrap
    
    Go1.3的标准库使用展位费 %w来实现wrap功能，但是不包含堆栈信息，功能上等同于withMessage

```