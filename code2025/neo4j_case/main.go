package main

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"log"
	"time"
)

/*

《neo4j 图数据库学习前置》

docker安装
```
docker run \
-d \
--restart=always \
--name neo4j \
-p 7474:7474 \
-p 7687:7687 \
-v neo4j:/data \
neo4j:4.4.5
```

控制台: http://10.0.40.3:7474/browser/
服务地址： 10.0.40.3:7687
账号: neo4j
密码: 123

go库教程: https://blog.csdn.net/lt326030434/article/details/124492583

*/

func main() {

	case2()
}

func case1() {

	neo4jURL := "bolt://10.0.40.3:7687"

	// 获取 neo4j driver 对象
	driver, err := neo4j.NewDriver(neo4jURL, neo4j.BasicAuth("neo4j", "123", ""))
	defer func(driver neo4j.Driver) {
		err = driver.Close()
		if err != nil {
			log.Println("neo4j close error:", err)
		}
	}(driver)

	if err != nil {
		log.Println("error connecting to neo4j:", err)
	}

	// 创建一个Session
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	// 开始一个事务
	tx, err := session.BeginTransaction()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	//// 创建一个节点
	//result, err := tx.Run("CREATE (n:Case1 {name: $value}) RETURN a", map[string]interface{}{
	//	"name": "BBBB",
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}

	//// case 1 使用给定的标签和属性创建节点。
	//// 批量创建
	//for i := 0; i < 10; i++ {
	//	result, err := tx.Run("CREATE (n:Case1 {id: $id, name: $name})", map[string]interface{}{
	//		"name": fmt.Sprintf("name-%d", i),
	//		"id":   i,
	//	})
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	// 处理结果
	//	if result.Next() {
	//		record := result.Record()
	//		node := record.GetByIndex(0).(neo4j.Node)
	//		fmt.Printf("Created Person node %d.\n", node.Id)
	//	} else {
	//		fmt.Println("No records found")
	//	}
	//}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("All done, sleeping for 5 seconds before exiting.")
	time.Sleep(5 * time.Second)

}

/*

CREATE (n:Label {name: $value})    使用给定的标签和属性创建节点。

CREATE (n:Label $map)    使用给定的标签和属性创建节点。

CREATE (n:Label)-[r:TYPE]->(m:Label)    根据给定的关系类型和方向创建关系；将变量r绑定到它。

CREATE (n:Label)-[:TYPE {name: $value}]->(m:Label)   使用给定的类型、方向和属性创建关系。

*/

// 给节点绑定关系
func case2() {
	neo4jURL := "bolt://10.0.40.3:7687"

	// 获取 neo4j driver 对象
	driver, err := neo4j.NewDriver(neo4jURL, neo4j.BasicAuth("neo4j", "123", ""))
	defer func(driver neo4j.Driver) {
		err = driver.Close()
		if err != nil {
			log.Println("neo4j close error:", err)
		}
	}(driver)

	// 创建一个会话
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	// 创建两个标签为Person的节点并绑定关系
	_, err = session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		// 创建节点
		result, err := transaction.Run(`
            CREATE (p1:Case2 {name: "Alice"})
            CREATE (p2:Case2 {name: "Bob"})
            RETURN p1, p2
        `, nil)
		if err != nil {
			return nil, err
		}

		// 获取创建的节点
		record, err := result.Single()
		if err != nil {
			return nil, err
		}

		// 获取节点
		p1, ok := record.Get("p1")
		if !ok {
			return nil, fmt.Errorf("not node")
		}
		p2, ok := record.Get("p2")
		if !ok {
			return nil, fmt.Errorf("not node")
		}

		// 创建关系
		_, err = transaction.Run(`
            MATCH (p1:Person), (p2:Person)
            WHERE p1.name = $name1 AND p2.name = $name2
            CREATE (p1)-[r:KNOWS]->(p2)
            RETURN r
        `, map[string]interface{}{
			"name1": "Alice",
			"name2": "Bob",
		})
		if err != nil {
			return nil, err
		}

		fmt.Printf("Created nodes %v and %v and relationship KNOWS between them\n", p1, p2)
		return nil, nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

/*

case3

-- 创建数据
CREATE (d:Dog{ id:1,name:"dog1",age:1 })
CREATE (c:Cat{ id:2,name:"cat1",age:2 })
-- 使用新节点创建关系
CREATE (d:Dog)-[r:Like]->(c:Cat)

-- 已知节点 创建带属性的关系 使用where
MATCH (d:Dog),(c:Cat)
WHERE d.name = 'dog1' and  c.name = 'cat1'
CREATE (d)-[r:love {info:"loveit"}]->(c)
RETURN r

-- 检索关系详情
MATCH (d)-[r:love]->(c)
RETURN d,r,c

-- 检索关系详情
MATCH (d)-[r:Like]->(c)
RETURN d,r,c



*/

/*

case4

-- 删除节点
DELETE <node-name-list>
-- 删除关系
DELETE <node1-name>,<node2-name>,<relationship-name>


*/

/*

case5

-- 创建多个标签多个属性
CREATE (m:Movie:Pic {id:'1', name:"test"})
-- 删除标签

-- 删除属性 并返回
MATCH (m:Movie:Pic)
REMOVE m.name
RETURN m

-- 删除标签
MATCH (m:Movie) REMOVE m:Pic

*/

/*

case6   set添加、更新属性

-- 添加属性
MATCH (book:Book)
SET book.title = 'superstar'
RETURN book

*/

/*

case7   ORDER BY排序

-- 默认升序
ORDER BY  <property-name-list>  [DESC]

-- 示例
MATCH (emp:Employee)
RETURN emp.empid,emp.name,emp.salary,emp.deptno
ORDER BY emp.name


*/

/*

case8  UNION合并

-- union合并
<MATCH Command1>
   UNION
<MATCH Command2>

-- 示例
MATCH (cc:CreditCard)
RETURN cc.id as id,cc.number as number,cc.name as name,
   cc.valid_from as valid_from,cc.valid_to as valid_to
UNION
MATCH (dc:DebitCard)
RETURN dc.id as id,dc.number as number,dc.name as name,
   dc.valid_from as valid_from,dc.valid_to as valid_to


*/

/*

case9  LIMIT和SKIP子句

-- 返回前两行
MATCH (emp:Employee)
RETURN emp
LIMIT 2

-- 跳过两行
MATCH (emp:Employee)
RETURN emp
SKIP 2

-- 实现分页，跳过1行然后返回2行
MATCH (emp:Employee)
RETURN emp
SKIP 1 LIMIT 2


*/

/*

case10  MERGE合并

MERGE = CREATE + MATCH
-- 语法
MERGE (<node-name>:<label-name>
{
   <Property1-name>:<Property1-Value>
   .....
   <Propertyn-name>:<Propertyn-Value>
})


*/

/*

case11

-- 过滤null
MATCH (e:Employee)
WHERE e.id IS NOT NULL
RETURN e.id,e.name,e.sal,e.deptno

MATCH (e:Employee)
WHERE e.id IS NULL
RETURN e.id,e.name,e.sal,e.deptno



*/

/*

case12 IN操作符

IN[<Collection-of-values>]
-- 示例
MATCH (e:Employee)
WHERE e.id IN [123,124]
RETURN e.id,e.name,e.sal,e.deptno


*/

/*

case13  INDEX索引

-- 创建索引语法
CREATE INDEX ON :<label_name> (<property_name>)
-- 删除索引
DROP INDEX ON :<label_name> (<property_name>)

-- 创建索引
CREATE INDEX ON :Customer (name)
-- 删除
DROP INDEX ON :Customer (name)


*/

/*
case14  UNIQUE唯一约束

-- 语法
CREATE CONSTRAINT ON (<label_name>) ASSERT <property_name> IS UNIQUE
DROP CONSTRAINT ON (<label_name>) ASSERT <property_name> IS UNIQUE

-- 示例
CREATE CONSTRAINT ON (cc:CreditCard) ASSERT cc.number IS UNIQUE
DROP CONSTRAINT ON (cc:CreditCard) ASSERT cc.number IS UNIQUE


*/

/*

case15  DISTINCT去重

-- 语法
MATCH (n:Movie) RETURN Distinct(n.name)

-- 示例
MATCH (n) RETURN distinct(n.name) LIMIT 25


*/

/*

参考资料
https://cloud.tencent.com/developer/article/2171687
https://cloud.tencent.com/developer/article/2110326

*/
