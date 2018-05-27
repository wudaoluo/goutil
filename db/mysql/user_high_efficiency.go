package mysql

import (
	"time"
	"fmt"
	"database/sql"
	"goutil/glog"
)

/*高效率多行数据操作 基本上都用到了事务 */

//insert
func (this *userService) High_efficiency_insert() error {
	start := time.Now()
	tx,err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		err := tx.Rollback()
		if err != sql.ErrTxDone && err != nil{
			glog.Error(err)
		}
	}()


	for i:=0;i<1000;i++ {
		//每次循环用的都是tx内部的连接，没有新建连接，所以效率高
		_,err:=tx.Exec("INSERT INTO user(user_name,weixin,email) values(?,?,?)","a","b","c")
		if err != nil {
			return err
		}
	}

	tx.Commit()
	end := time.Now()
	fmt.Println("高效率insert ",end.Sub(start).Seconds())
	return nil
}



//update
func (this *userService) High_efficiency_update() error {
	start := time.Now()

	tx,err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		err := tx.Rollback()
		if err != sql.ErrTxDone && err != nil{
			glog.Error(err)
		}
	}()

	for i:=0;i<1000;i++ {
		//每次循环用的都是tx内部的连接，没有新建连接，所以效率高
		_,err:=tx.Exec("UPDATE user set user_name=? where id = ?","c",i)
		if err != nil {
			return err
		}
	}

	tx.Commit()
	end := time.Now()
	fmt.Println("高效率update ",end.Sub(start).Seconds())
	return nil
}

//delete
func (this *userService) High_efficiency_delete() error {
	start := time.Now()

	tx,err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		err := tx.Rollback()
		if err != sql.ErrTxDone && err != nil{
			glog.Error(err)
		}
	}()


	for i:=0;i<1000;i++ {
		//每次循环用的都是tx内部的连接，没有新建连接，所以效率高
		_,err:=tx.Exec("DELETE FROM user WHERE id=?",i)
		if err != nil {
			return err
		}
	}

	tx.Commit()
	end := time.Now()
	fmt.Println("高效率delete ",end.Sub(start).Seconds())
	return nil
}


//select
func (this *userService) High_efficiency_select() ([]User,error) {
	start := time.Now()
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

	end := time.Now()
	fmt.Println("高效率select ",end.Sub(start).Seconds())

	return list,nil
}

/*

    //方式2 query
    start = time.Now()
    stm,_ := db.Prepare("SELECT uid,username FROM USER")
    defer stm.Close()
    rows,_ = stm.Query()
    defer rows.Close()
    for rows.Next(){
         var name string
         var id int
        if err := rows.Scan(&id,&name); err != nil {
            log.Fatal(err)
        }
       // fmt.Printf("name:%s ,id:is %d\n", name, id)
    }
    end = time.Now()
    fmt.Println("方式2 query total time:",end.Sub(start).Seconds())
*/