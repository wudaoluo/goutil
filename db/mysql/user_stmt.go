package mysql

import (
	"log"
	"database/sql"
	"fmt"
)

/*
	这种方式适合多条数据处理 ，底层要比较复杂
	尽量不要使用
*/


func (this *userService) Stmt_select_one() (*User,error) {
	stmt, err := db.Prepare("SELECT * FROM user where user_name = ? limit 1")
	if err != nil {
		log.Fatalln(err)
	}
	defer stmt.Close()

	row :=  stmt.QueryRow("xys")

	u := &User{}
	err = row.Scan(&u.Id,&u.UserName,&u.Weixin,&u.Email,&u.Status)
	if err != nil{
		//只有当查询的结果为空的时候，会触发一个sql.ErrNoRows错误
		if err == sql.ErrNoRows{
			return nil,fmt.Errorf("user not found")
		}else {
			return nil,err
		}
	}

	return u,nil
}


func (this *userService) Stmt_select() ([]User,error){
	stmt, err := db.Prepare("SELECT * FROM user where user_name = ?")
	if err != nil {
		log.Fatalln(err)
	}
	defer stmt.Close()

	rows, err :=  stmt.Query("xys")
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




func (this *userService) Stmt_insert() (int64,error) {
	stmt, err := db.Prepare("INSERT INTO user (user_name) VALUES (?)")

	if err != nil {
		return 0,err
	}

	defer stmt.Close()

	rows, err := stmt.Exec("xys")
	if err != nil {
		return  0,err
	}

	last_insert_id ,err := rows.LastInsertId()
	if err != nil {
		return 0,err
	}
	return last_insert_id,err
}


func (this *userService) stmt_update() (int64,error){
	stmt, err := db.Prepare("UPDATE user set weixin = ?")

	if err != nil {
		return 0,err
	}

	defer stmt.Close()

	rows, err := stmt.Exec("123wexin")
	if err != nil {
		return  0,err
	}

	rows_affected ,err := rows.RowsAffected()
	if err != nil {
		return 0,nil
	}
	return rows_affected,err

}

func (this *userService) stmt_delete() (int64,error){
	stmt, err := db.Prepare("DELETE FROM user where id = ?")

	if err != nil {
		return 0,err
	}

	defer stmt.Close()

	rows, err := stmt.Exec(23)
	if err != nil {
		return  0,err
	}

	rows_affected ,err := rows.RowsAffected()
	if err != nil {
		return 0,err
	}
	return rows_affected,err

}