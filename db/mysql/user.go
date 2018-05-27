package mysql

import (
	"time"
	"database/sql"
	"fmt"
	"goutil/glog"
)


//在创建表的时候要给默认值，null值对整个系统都增加了复杂度
//------增删改查的例子----------/


type User struct {
	Id 			int
	UserName 	string
	Weixin 		string
	Email 		string
	Status 		int            // 状态，1正常 0禁用
	CreateTime  time.Time      // 创建时间
}


type userService struct {}


func (this *userService) table() string {
	return tableName("user")
}



/*
在实际应用中，与数据库交互，往往写的sql语句还带有参数，这类sql可以称之为prepare语句。
prepare语句有很多好处，可以防止sql注入，可以批量执行等。
但是prepare的连接管理有其自己的机制。


对于写，即插入更新和删除。这类操作与query不太一样，写的操作只关系是否写成功了。
database/sql提供了Exec方法用于执行写的操作。
我们也见识到了，Eexc返回一个sql.Result类型，
它有两个方法LastInsertId和RowsAffected。LastInsertId返回是一个数据库自增的id，这是一个int64类型的值。
Exec执行完毕之后，连接会立即释放回到连接池中，因此不需要像query那样再手动调用row的close方法。


尽量避免使用prepare方式
*/



//查询单行
func (this *userService) GetbyId(id int) (*User,error) {
	row := db.QueryRow("select id,user_name,weixin,email,status from user where id = ?",id)
	u := &User{}
	err := row.Scan(&u.Id,&u.UserName,&u.Weixin,&u.Email,&u.Status)
	if err != nil{
		//只有当查询的结果为空的时候，会触发一个sql.ErrNoRows错误
		if err == sql.ErrNoRows{
			return nil,fmt.Errorf("user not found id: ",id)
		}else {
			return nil,err
		}
	}

	return u,nil
}


//查询多行
func (this *userService) GetUserList() ([]User,error) {
	rows, err := db.Query("select id,user_name,weixin,email,status from user")
	if err != nil {
		return nil ,err
	}

	list := make([]User,0)

	for rows.Next() {
		u := new(User)

		err = rows.Scan(&u.Id,&u.UserName,&u.Weixin,&u.Email,&u.Status)
		if err != nil {
			rows.Close()
			return nil,err
		}

		list = append(list,*u)
	}

	rows.Close()

	//为了检查是否是迭代正常退出还是异常退出，需要检查rows.Err
	if rows.Err() != nil {
		return nil,err
	}

	return list,nil
}


//查询多行 只需要输入sql语句  SELECT * FROM user
func GetResult(db *sql.DB, sql string) ([]map[string]interface{},error) {
	rows, err := db.Query(sql)

	if err != nil {
		return nil,err
	}

	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil,err
	}

	count := len(cols)
	vals := make([][]byte, count)
	scans := make([]interface{},count)

	var list []map[string]interface{}

	for i := 0; i < count; i++ {
		//ptr[i] = &values[i]
		scans[i] = &vals[i]
	}

	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil{
			return nil,err
		}

		row := make(map[string]interface{})
		for k, v := range vals{
			key := cols[k]
			row[key] = string(v)
		}
		list = append(list,row)



	}

	return list,nil

}



//insert
func (this *userService) Add(UserName string) error {
	_, err := db.Exec("INSERT INTO user(user_name) values(?)",UserName)
	return err

	// n,err := row.LastInsertId() 返回插入的id
	// n,err := row.RowsAffected() 返回受影响的行数
}


//update
func (this *userService) Update(id int) error {
	_, err := db.Exec("DELETE from user where id ?",id)
	return err
}


//delete
func (this *userService) Delete(id int) error {
	_, err := db.Exec("DELETE from user where id ?",id)
	return err
}


//事务处理
/*
因为事务是单个连接，因此任何事务处理过程的出现了异常，都需要使用rollback，一方面是为了保证数据完整一致性，另一方面是释放事务绑定的连接。
*/
func (this *userService) shiwu() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		err := tx.Rollback()
		if err != sql.ErrTxDone && err != nil{
			glog.Error(err)
		}
	}()

	_, err = tx.Exec("UPDATE user SET Status=1 WHERE user_name='vanyarpy'")
	if err != nil {
		return err
	}


	rs, err := tx.Exec("UPDATE user SET Status=1 WHERE user_name='noldorpy'")
	if err != nil {
		return err
	}
	rowAffected, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println(rowAffected)


	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}







