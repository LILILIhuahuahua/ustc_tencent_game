package collision

import (
	"container/list"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
)

// QuadTree	四叉树，用于在碰撞检测时快速返回目标物体的邻近物体
type QuadTree struct {
	maxObjectNum int32	//区域内可容纳的最大物体数量
	maxLevelNum  int32	//四叉树最大层数
	curLevel int32	//四叉树当前层数
	bounds	*Rectangle	//四叉树边界数据
	objects *list.List	//四叉树中物体集合
	childs []*QuadTree	//四叉树子节点集合
}

// @title    NewQuadTree
// @description   四叉树构造函数
// @param     curLevel        int32         "四叉树当前层数"
// @param     bounds        *Rectangle         "四叉树区域边界"
// @return    tree        *QuadTree         "四叉树指针"
func NewQuadTree(curLevel int32, bounds *Rectangle) *QuadTree{
	tree := &QuadTree{}
	tree.maxObjectNum = configs.MaxObjectNum
	tree.maxLevelNum = configs.MaxLevelNum
	tree.curLevel = curLevel
	tree.bounds = bounds
	tree.objects = list.New()
	tree.childs = make([]*QuadTree, 4)
	return tree
}

// @title    Split
// @description   四叉树节点分裂(由于当前节点中物体数量大于区域内可容纳的最大物体数量)
func (tree *QuadTree)Split()  {
	halfX := (tree.bounds.MaxX - tree.bounds.MinX)/2
	halfY := (tree.bounds.MaxY - tree.bounds.MinY)/2
	originX := tree.bounds.GetX()
	originY := tree.bounds.GetY()
	child0 := NewQuadTree(tree.curLevel+1, NewRectangleByBounds(originX, originY, originX+halfX, originY+halfY))
	child1 := NewQuadTree(tree.curLevel+1, NewRectangleByBounds(originX-halfX, originY, originX, originY+halfY))
	child2 := NewQuadTree(tree.curLevel+1, NewRectangleByBounds(originX-halfX, originY-halfY, originX, originY))
	child3 := NewQuadTree(tree.curLevel+1, NewRectangleByBounds(originX, originY-halfY, originX+halfX, originY))
	tree.childs = append(tree.childs, child0)
	tree.childs = append(tree.childs, child1)
	tree.childs = append(tree.childs, child2)
	tree.childs = append(tree.childs, child3)
}

// @title    GetDistrictIndex
// @description   查询物体所在四叉树区域编号（0：第一象限 1：第二象限 2：第三象限 3：第四象限）
// @param     obj        *Rectangle         "物体"
// @return    index        int32        "物体所在区域编号"
func (tree *QuadTree)GetDistrictIndex(obj *Rectangle) int32 {
	var index int32 = -1 //-1代表属于本节点，例如当物体位置恰好在分界线上时该物体不属于任何一个象限
	if obj.GetY() > tree.bounds.GetY() && obj.GetX() > tree.bounds.GetX() {
		index = 0
	}
	if obj.GetY() > tree.bounds.GetY() && obj.GetX() < tree.bounds.GetX() {
		index = 1
	}
	if obj.GetY() < tree.bounds.GetY() && obj.GetX() < tree.bounds.GetX() {
		index = 2
	}
	if obj.GetY() < tree.bounds.GetY() && obj.GetX() > tree.bounds.GetX() {
		index = 3
	}
	return index
}

// @title    InsertObj
// @description   向四叉树中插入物体
// @param     obj        *Rectangle         "物体"
func (tree *QuadTree)InsertObj(obj *Rectangle)  {
	if len(tree.childs) > 0 {
		index := tree.GetDistrictIndex(obj)
		if -1 != index {
			tree.childs[index].InsertObj(obj)
			return
		}
	}
	//tree.objects = append(tree.objects, obj)
	tree.objects.PushBack(obj)
	if tree.objects.Len() > int(tree.maxObjectNum) {
		tree.Split()	//分裂下一层节点
		var next *list.Element
		for e:= tree.objects.Front(); e!=nil; e=next {	//将本层中的物体移动至下一层
			next = e.Next()
			eIndex := tree.GetDistrictIndex(obj)
			if -1 != eIndex {
				tree.objects.Remove(e)
				tree.childs[eIndex].InsertObj(e.Value.(*Rectangle))
			}
		}
	}
}

// @title    GetObjsInSameDistrict
// @description   查找四叉树中目标物体的同区物体集合
// @param     obj        *Rectangle         "目标物体"
// @return    objs        []*Rectangle      "同区物体集合"
func (tree *QuadTree)GetObjsInSameDistrict(obj *Rectangle) []*Rectangle {
	objs := []*Rectangle{}
	index := tree.GetDistrictIndex(obj)
	if -1 != index && len(tree.childs) > 0{
		return tree.childs[index].GetObjsInSameDistrict(obj)
	}
	for e:= tree.objects.Front();e!=nil;e=e.Next() {
		if e.Value.(*Rectangle) != obj { //去除自己
			objs = append(objs, e.Value.(*Rectangle))
		}
	}
	return objs
}

// @title    Clear
// @description   清除整个四叉树中的所有物体
func (tree *QuadTree)Clear()  {
	var next *list.Element
	for e:= tree.objects.Front(); e!=nil; e=next {	//将本层中的物体移动至下一层
		next = e.Next()
		tree.objects.Remove(e)
	}
	for i,child := range tree.childs {
		if nil != child {
			child.Clear()
			tree.childs[i] = nil
		}
	}
}





