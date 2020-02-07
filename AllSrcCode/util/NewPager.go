package util

import (
	"bytes"
	"fmt"
	"math"
	"strings"
)

type Pager struct{
	Page int
	Totalnum int
	Pagesize int
	urlpath string
	urlquery string
	nopath bool
}

func NewPager(page, totalnum,pagesize int, url string, nopath ...bool) *Pager{
	p:=Pager{
		Page:page,
		Totalnum:totalnum,
		Pagesize:pagesize,
	}
	//解析url
	arr:=strings.Split(url,"?")
	p.urlpath=arr[0]
	if len(arr)>1{ //url中包含了?
		p.urlquery="?"+arr[1]
	}else{
		p.urlquery=""
	}

	if len(nopath)>0{
		p.nopath=nopath[0]
	}else{
		p.nopath=false
	}
	//要求返回一个指针
	return &p
}
func (this *Pager) url(page int) string {
	if this.nopath { //不使用目录形式
		if this.urlquery != "" {
			return fmt.Sprintf("%s%s&page=%d", this.urlpath, this.urlquery, page)
		} else {
			return fmt.Sprintf("%s?page=%d", this.urlpath, page)
		}
	} else {
		return fmt.Sprintf("%s/page/%d%s", this.urlpath, page, this.urlquery)
	}
}
func (this *Pager) ToString() string{
	if this.Totalnum<=this.Pagesize{
		return ""
	}
	var buf bytes.Buffer
	var from, to, linknum, offset, totalpage int

	offset = 5
	linknum = 10

	//计算页数
	totalpage = int(math.Ceil(float64(this.Totalnum) / float64(this.Pagesize)))
	if totalpage < linknum {
		from = 1
		to = totalpage
	} else { //如果页数大于10，只显示前后5页
		from = this.Page - offset
		to = from + linknum
		if from < 1 {
			from = 1
			to = from + linknum - 1
		} else if to > totalpage {
			to = totalpage
			from = totalpage - linknum + 1
		}
	}

	//将页数信息更新到前端
	buf.WriteString("<div class=\"pagination\"><ul>")
	if this.Page > 1 {
		buf.WriteString(fmt.Sprintf("<li><a href=\"%s\">&laquo;</a></li>", this.url(this.Page-1)))
	} else {
		buf.WriteString("<li class=\"disabled\"><span>&laquo;</span></li>")
	}

	if this.Page > linknum {
		buf.WriteString(fmt.Sprintf("<li><a href=\"%s\">1...</a></li>", this.url(1)))
	}

	for i := from; i <= to; i++ {
		if i == this.Page {
			buf.WriteString(fmt.Sprintf("<li class=\"active\"><span>%d</span></li>", i))
		} else {
			buf.WriteString(fmt.Sprintf("<li><a href=\"%s\">%d</a></li>", this.url(i), i))
		}
	}

	if totalpage > to {
		buf.WriteString(fmt.Sprintf("<li><a href=\"%s\">...%d</a></li>", this.url(totalpage), totalpage))
	}

	if this.Page < totalpage {
		buf.WriteString(fmt.Sprintf("<li><a href=\"%s\">&raquo;</a></li>", this.url(this.Page+1)))
	} else {
		buf.WriteString(fmt.Sprintf("<li class=\"disabled\"><span>&raquo;</span></li>"))
	}
	buf.WriteString("</ul></div>")

	return buf.String()
}