# []()资料

比较好的Gorm中文文档 https\://jasperxu.com/gorm-zh 本文基于该资料进行整理，汇总最基本的Gorm入门使用内容，陆续补充。

# []()安装

```bash
go get -u github.com/jinzhu/gorm
```

# []()数据库配置

```go
//数据库配置信息
func options() {
	// 全局禁用表名复数
	// 如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响
	DB.SingularTable(true)
	//自动迁移模式将保持更新到最新。
	//警告：自动迁移仅仅会创建表，缺少列和索引，并且不会改变现有列的类型或删除未使用的列以保护数据。
	DB.AutoMigrate(&User{})
	//打印日志，默认是false
	DB.LogMode(true)
	//连接池
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
}
```

# []()数据库连接

## []()sqlite3

使用`sqlite`数据库来快速连接

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	//这里是gorm封装的数据库驱动
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func initSqlite3Db() {
	// 初始化
	var err error
	DB, err = gorm.Open("sqlite3", "test.db")
	// 检查错误
	if err != nil {
		panic(err)
	}
}
```

## []()mysql

```go
func initMysqlDb() {
	// 初始化
	var err error
	DB, err = gorm.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/testdb?charset=utf8")
	// 检查错误
	if err != nil {
		panic(err)
	}
}
```

# []()模型定义

## []()tag:gorm

通过`tag`标记`gorm`来创建在数据库中字段的约束和属性配置

```go
type User struct {
	Id         int64  `gorm:"primary_key;AUTO_INCREMENT"` //设置为primary_key主键，AUTO_INCREMENT自增
	UserId     string  `gorm:"index:idx_user_id;not null`     //设置普通索引，索引名称为idx_user_id,not null不能为空
	UserName   string `gorm:"type:varchar(64);not null"`  //设置为varchar类型，长度为64，not null不能为空
	Age        int
	Phone      int       `gorm:"unique_index:uk_phone"` //设置唯一索引，索引名为uk_phone
	CreateTime time.Time `gorm:"not null"`
	UpdateTime time.Time `gorm:"not null"`
	UserRemark string    `gorm:"column:remark;default:'默认'"`
	IgnoreMe   int       `gorm:"-"` // 忽略这个字段
}
```

| gorm定义                    | 数据库字段约束                 |
| ------------------------- | ----------------------- |
| `gorm:"primary_key"`      | 字段设置为主键                 |
| `gorm:"AUTO_INCREMENT"`   | 字段设置为自增                 |
| `gorm:"size:20`           | 字段长度设置为20               |
| `gorm:"index:idx_user_id` | 字段设置普通索引，名称为idx_user_id |
| `gorm:"not null`          | 设置字段为非空                 |
| `gorm:"type:varchar(64)"` | 设置字段为varchar类型，长度为64    |
| `gorm:"column:remark"`    | 设置数据库字段名为remark         |
| `gorm:"-"`                | 忽略此字段，不在表中创建该字段         |
| `gorm:"default:'默认'"`     | 设置字段的默认值                |

> 注意：不一定是最全的示例，陆续补充

## []()表名

在`gorm`中，表名创建是以复数存在的，当创建表之前设置以下选项可以禁用表名复数

```go
	// 全局禁用表名复数
	// 如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响
	DB.SingularTable(true)
```

或者在每个`model`结构体中声明`TableName`方法，若存在该方法`gorm`会**优先**采用这种方式来配置表名

```go
type User struct {} // 默认表名是`users`

// 设置User的表名为`user`
func (User) TableName() string {
  return "user"
}
```

# []()基础操作

## []()新增

### []()NewRecord主键检查 & Create

```go
//Create
func CreateUser() {
	user := &User{
		UserId:     "USER_001",
		UserName:   "zhangsan",
		Age:        10,
		Phone:      19901020305,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		UserRemark: "备注",
	}
	//1、NewRecord会检查主键是否存在
	boolBefore := DB.NewRecord(user)
	//主键存在返回true；否则返回false	这里返回true，因为主键ID是自动生成这里ID是空的不存在
	fmt.Println("boolBefore => ", boolBefore)

	//2、Create
	err := DB.Create(user).Error
	if err != nil {
		fmt.Printf("CreateUser error %v \n", err)
	}
	//这里会返回主键ID
	fmt.Println("user => ", toJson(user))

	//3、NewRecord会检查主键是否存在
	boolAfter := DB.NewRecord(user)
	//主键存在返回true；否则返回false  这里返回false，因为主键ID已经存在了
	fmt.Println("boolAfter => ", boolAfter)
	//这里有完整数据返回，包含主键ID
	fmt.Println("user => ", toJson(user))
}
```

打印日志如下：

```bash
boolBefore =>  true
user =>  {"Id":1,"UserId":"USER_001","UserName":"zhangsan","Age":10,"Phone":19901020301,"CreateTime":"2021-06-01T13:33:02.237956+08:00","UpdateTime":"2021-06-01T13:33:02.237956+08:00","UserRemark":"默认","IgnoreMe":0}
boolAfter =>  false
user =>  {"Id":1,"UserId":"USER_001","UserName":"zhangsan","Age":10,"Phone":19901020301,"CreateTime":"2021-06-01T13:33:02.237956+08:00","UpdateTime":"2021-06-01T13:33:02.237956+08:00","UserRemark":"默认","IgnoreMe":0}
```

表数据如下：

```sql
sqlite> select * from user;
1|USER_001|zhangsan|10|19901020301|2021-06-01 13:33:02.237956+08:00|2021-06-01 13:33:02.237956+08:00|默认
```

## []()查询

### []()First:查询第一条记录

```go
func FindFirstRecord()  {
	result := &User{}
	// SELECT * FROM user ORDER BY id LIMIT 1;
	err:=DB.First(result).Error
	
	if err != nil {
		fmt.Printf("FindFirstRecord error %v",err)
		return
	}
	fmt.Printf(toJson(result))
}
```

日志输出如下：

```bash
{"Id":1,"UserId":"USER_001","UserName":"zhangsan","Age":10,"Phone":19901020301,"CreateTime":"2021-06-01T13:33:02.237956+08:00","UpdateTime":"2021-06-01T13:33:02.237956+08:00","UserRemark":"默认","IgnoreMe":0}
```

### []()Last:查询最后一条记录

```go
func FindLastRecord()  {
	result := &User{}
	// SELECT * FROM users ORDER BY id DESC LIMIT 1
	err:=DB.Last(result).Error

	if err != nil {
		fmt.Printf("FindLastRecord error %v",err)
	    return
	}
	fmt.Printf(toJson(result))
}
```

日志输出如下：

```bash
{"Id":2,"UserId":"USER_001","UserName":"zhangsan","Age":10,"Phone":19901020302,"CreateTime":"2021-06-01T13:57:56.807786+08:00","UpdateTime":"2021-06-01T13:57:56.807786+08:00","UserRemark":"默认","IgnoreMe":0}
```

### []()First(… , pk):根据主键查询记录

```go
func FindRecordByPK()  {
	result := &User{}
	// SELECT * FROM users WHERE id = 2
	err:=DB.First(&User{},2).Find(result).Error

	if err != nil {
		fmt.Printf("FindRecordByPK error %v",err)
		return
	}
	fmt.Printf(toJson(result))
}
```

日志输出如下：

```bash
{"Id":2,"UserId":"USER_001","UserName":"zhangsan","Age":10,"Phone":19901020302,"CreateTime":"2021-06-01T13:57:56.807786+08:00","UpdateTime":"2021-06-01T13:57:56.807786+08:00","UserRemark":"默认","IgnoreMe":0}
```

### []()Where(…) 条件查询条件

关于条件查询相对来说比较简单，这里不赘述可以参考 https\://jasperxu.com/gorm-zh/crud.html#q

### []()FirstOrInit() & Attrs() & Assign()

获取第一个匹配的记录，或者使用给定的条件初始化一个新的记录返回。

```go
func FirstOrInit() {
	user := &User{
		UserId:     "USER_001",
		UserName:   "User",
		Age:        10,
		Phone:      19901020305,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	// 查询 SELECT * FROM users WHERE user_name="guanjian" 该记录存在
	err := DB.Attrs(&User{UserName: "attrs"}).
		// 这里无论是否找到都强制赋值
		Assign(&User{UserRemark: "强制修改"}).
		// 没有这条记录则填补,初始化的填充逻辑是 Attrs > Where (没有Attrs则使用Where条件)
		// 如果有则返回查到的第一条匹配的记录
		FirstOrInit(user, User{UserName: "guanjian"}).Error

	if err != nil {
		fmt.Printf("FirstOrInit error %v", err)
		return
	}
	fmt.Printf(toJson(user))
}
```

> `FirstOrInit`方法如果查不到结果**不会落库**，只是会初始化`FirstOrInit(obj)`中的`obj`对象的值，且将Where条件中的值一并初始化出来

### []()FirstOrCreate() & Attrs() & Assign()

和`FirstOrInit`类似，不同的是它会进行**落库**操作。获取第一个匹配的记录，或者使用给定的条件进行落库新建。

```go
func FirstOrCreate() {
	user := &User{
		UserId:     "USER_001",
		UserName:   "guanjian",
		Age:        10,
		Phone:      19901020310,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	// 查询 SELECT * FROM users WHERE id = 10
	// 没有这条记录则按照user来初始化user进行填补
	// 如果有则返回查到的第一条匹配的记录
	err := DB.FirstOrCreate(user, User{Id: 10}).Error

	if err != nil {
		fmt.Printf("FirstOrInit error %v", err)
		return
	}
	fmt.Printf(toJson(user))
}
```

> `FirstOrCreate`方法如果查不到结果**会落库**，也会初始化`FirstOrCreate(obj)`中的`obj`对象的值，且将Where条件中的值一并初始化出来

### []()Select查询字段定义

这里只查询 `user_name,age`字段返回，其余字段不会查询到值返回

```go
func Select() {
	user := &User{}
	err := DB.Select("user_name,age").Find(user).Error

	if err != nil {
		fmt.Printf("Select error %v", err)
		return
	}
	fmt.Println("只查询user_name,age字段", toJson(user))
}
```

### []()Scan

```go
func Scan() {
	user := &User{}
	err := DB.Table("user").Select([]string{"user_name","age"}).Scan(user).Error
	if err != nil {
		fmt.Printf("Scan error %v", err)
		return
	}
	fmt.Println("user => ", toJson(user))
}
```

### []()Scopes 动态条件添加

实现方法形如`func(db *gorm.DB) * gorm.DB`则可以配合方法`Scopes`进行动态条件的添加，编码更清晰友好

```go
//这里匹配create_time查询条件
func WithCreateTimeLessThanNow(db *gorm.DB) *gorm.DB {
	return db.Where("create_time < ?", time.Now())
}

//这里匹配age查询条件
func WithAgeGreaterThan0(db *gorm.DB) *gorm.DB {
	return db.Where("age > ?", 0)
}

func Scopes() {
	user := &User{}
	//查询条件动态拼接
	err := DB.Scopes(WithAgeGreaterThan0,WithCreateTimeLessThanNow).Find(user).Error
	if err != nil {
		fmt.Printf("Scopes error %v", err)
		return
	}
	fmt.Println("user => ", toJson(user))
}
```

## []()修改

### []()Save

保存已存在的数据，触发更新逻辑（update）；保存不存在的数据，直接插入（insert）

```go
//保存存在的数据，触发更新逻辑
func SaveExist() {
	user := &User{
		Id: 1,
	}
	//查询id=9999的数据，是不存在的
	user.UserName = "saveNewUserName"
	//按照主键做Update操作
	err := DB.Save(user).Error
	if err != nil {
		fmt.Printf("Save error %v", err)
		return
	}
	fmt.Println("save user => ", toJson(user))
}
//保存不存在的数据，直接插入
func SaveNoExist() {
	user := &User{
		Id: 9999,
		Phone: 19902029492,
	}
	//查询id=9999的数据，是不存在的
	user.UserName = "saveNewUserName"
	//按照主键做Update操作
	err := DB.Save(user).Error
	if err != nil {
		fmt.Printf("Save error %v", err)
		return
	}
	fmt.Println("save user => ", toJson(user))
}
```

### []()Update & Updates & Omit

```go
func Update() {
	user := &User{
		Id: 9999,
	}

	err := DB.Model(user).Update("user_name", "this is 9999").Error
	if err != nil {
		fmt.Printf("Update error %v", err)
		return
	}
}
```

更新字段效果如下：\
![在这里插入图片描述](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/fa32bb316aae4eb2bfda3e21f6ee77e3~tplv-k3u1fbpfcp-zoom-1.image)\
使用`Updates`方法可以批量更新多个字段

```go
func Updates() {
	//查询字段
	whereUser := &User{
		Id: 9999,
	}
	//批量更新字段
	updateUser := &User{
		UserName: "This is Updates",
		Age: 9292,
		UserRemark: "Updates",
	}

	err := DB.Model(whereUser).Updates(updateUser).Error
	if err != nil {
		fmt.Printf("Updates error %v", err)
		return
	}
}
```

批量更新字段效果如下：\
![在这里插入图片描述](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/f36bc1f05b3f4caf8b34c9473564aad8~tplv-k3u1fbpfcp-zoom-1.image)\
使用`Omit`方法可以忽略该字段的更新

```go
func UpdateOmit() {
	//查询字段
	whereUser := &User{
		Id: 1,
	}
	//批量更新字段
	updateUser := &User{
		UserName:   "This is Updates",
		Age:        6666,
		UserRemark: "1 -> Update",
	}

	err := DB.Model(whereUser).
		//忽略user_name字段的更新
		Omit("user_name").
		Updates(updateUser).Error
	if err != nil {
		fmt.Printf("Updates error %v", err)
		return
	}
}
```

![在这里插入图片描述](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/27c8dd90b8874a9f98b5191cca1c3f9e~tplv-k3u1fbpfcp-zoom-1.image)

> 以上更新操作将执行模型的`BeforeUpdate, AfterUpdate`方法，更新其UpdatedAt时间戳，在更新时保存它的Associations，如果不想调用它们，可以使用UpdateColumn, UpdateColumns\
> \
> **非零值更新问题** 为了避免误操作或脏数据对全表进行更新，需要配置`WHERE`必填条件进行限制

使用`RowsAffected`可以返回当前更新的行数

```go
func RowsAffected() {
	//查询字段
	whereUser := &User{
		Id: 9999,
	}
	//批量更新字段
	updateUser := &User{
		UserName:   "This is Updates",
		Age:        9292,
		UserRemark: "Updates",
	}
	//RowsAffected返回当前更新的行数
	count := DB.Model(whereUser).Updates(updateUser).RowsAffected
	fmt.Printf("Updates count %v", count)
}
```

## []()删除

一般使用数据库逻辑删除，这部分可参考 https\://jasperxu.com/gorm-zh/crud.html#d

# []()回调使用（Callback）

在创建，更新，查询，删除时将被调用，如果任何回调返回错误，gorm将停止未来操作并回滚所有更改。

## []()Callback默认定义

在`gorm`源码中可以看到默认的`callback`注册，这里着重关注`callback_create.go、callback_delete.go、callback_query.go、callback_update.go`即可，其实整个`gorm`的执行流程都是基于`callback`来实现整个执行流程的。\
![在这里插入图片描述](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/023499c28a7c453496ef2b774c3d04dc~tplv-k3u1fbpfcp-zoom-1.image)\
下面举例`callback_create.go`文件的`init()`方法，按照从上到下的顺序，注册了`callback`函数，每个`callbackName`（如：gorm:create）都可以作为回调的插入点使用，配合`Before、After`可以实现前置插入、后置插入等，其他`callback`类似这里不再赘述

```go
// Define callbacks for creating
func init() {
	DefaultCallback.Create().Register("gorm:begin_transaction", beginTransactionCallback)
	DefaultCallback.Create().Register("gorm:before_create", beforeCreateCallback)
	DefaultCallback.Create().Register("gorm:save_before_associations", saveBeforeAssociationsCallback)
	DefaultCallback.Create().Register("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	DefaultCallback.Create().Register("gorm:create", createCallback)
	DefaultCallback.Create().Register("gorm:force_reload_after_create", forceReloadAfterCreateCallback)
	DefaultCallback.Create().Register("gorm:save_after_associations", saveAfterAssociationsCallback)
	DefaultCallback.Create().Register("gorm:after_create", afterCreateCallback)
	DefaultCallback.Create().Register("gorm:commit_or_rollback_transaction", commitOrRollbackTransactionCallback)
}
```

## []()Callback事务保证

查看`gorm`源码，诸如CRUD的方法如`Create()、Update()、Find()、Delete()`等内部都进行了callback方法等调用，callback内部代码如下：

```go
func (scope *Scope) callCallbacks(funcs []*func(s *Scope)) *Scope {
	//这里进行了异常捕获
	defer func() {
		if err := recover(); err != nil {
			if db, ok := scope.db.db.(sqlTx); ok {
				//异常则对事务进行回滚
				db.Rollback()
			}
			panic(err)
		}
	}()
	for _, f := range funcs {
		(*f)(scope)
		if scope.skipLeft {
			break
		}
	}
	return scope
}
```

下面来看一个`callback`中产生panic异常，数据库回滚的示例，如下：

```go
func Callback() {
	//注册在gorm:create节点之后执行函数
	DB.Callback().Create().
		//在CreateUser()函数方法，即DB执行之后再执行
		After("gorm:create").
		Register("gorm:AfterCreate", func(scope *gorm.Scope) {
		fmt.Println("AfterCreate")
		//这里抛出异常，会被callCallbacks的recover方法捕获
		panic("AfterCreate error")
	})
	//由于上面callback已经panic，这里会回滚
	CreateUser()
}
```

# []()异常处理

可以对`RecordNotFound()、Error、GetErrors()`来处理记录不存在、异常、多个异常等情况

```go
func Error() {
	//RecordNotFound
	bool := DB.First(&User{Id: 9999}).RecordNotFound()
	fmt.Println("RecordNotFound => ", bool)

	//Error
	err := DB.Create(&User{}).Error
	fmt.Println("err => ", err)

	errs := DB.Create(&User{}).GetErrors()
	fmt.Println("errs => ", errs)
}
```

# []()事务

注意，一旦你在一个事务中，使用`tx`作为数据库句柄，**如果使用`gorm.DB`操作，是不受事务管控的，这里需要注意！**

```go
func Tx() {
	//开启事务
	tx := DB.Begin()

	//执行方法1  会执行成功
	err1 := tx.Create(&User{
		UserId:     "USER_002",
		UserName:   "zhangsan",
		Age:        0,
		Phone:      17701020304,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		UserRemark: "",
		IgnoreMe:   0,
	}).Error
	if err1 != nil {
		tx.Rollback()
		fmt.Println("=== Create 1 error ===")
	}
	//执行方法2 会执行失败，有Phone唯一主键约束
	err2 := tx.Create(&User{
		UserId:     "USER_002",
		UserName:   "zhangsan",
		Age:        0,
		Phone:      17701020304,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		UserRemark: "",
		IgnoreMe:   0,
	}).Error
	if err2 != nil {
		tx.Rollback()
		fmt.Println("=== Create 2 error ===")
	}
	tx.Commit()
}
```

# []()日志

```go
// 启用Logger，显示详细日志
db.LogMode(true)

// 禁用日志记录器，不显示任何日志
db.LogMode(false)

// 调试单个操作，显示此操作的详细日志
db.Debug().Where("name = ?", "jinzhu").First(&User{})
```

# []()参考

https\://jasperxu.com/gorm-zh/
