// Copyright 2013 wetalk authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package post

import (
	"github.com/beego/wetalk/modules/models"
	"github.com/beego/wetalk/modules/post"
	"github.com/beego/wetalk/routers/base"
	"github.com/beego/wetalk/setting"
)

//Post List Router

type PostListRouter struct {
	base.BaseRouter
}

//Get all the categories
func (this *PostListRouter) setCategories(cats *[]models.Category) {
	//@see modules/post/topic_util.go
	post.ListCategories(cats)
	this.Data["Categories"] = *cats
}

//Get all the topics of the category
func (this *PostListRouter) setTopicsOfCategory(topics *[]models.Topic, category *models.Category) {
	//@see modules/post/topic_util.go
	post.ListTopicsOfCategory(topics, category)
	this.Data["TopicsOfCategory"] = *topics
}

//Get the home page
func (this *PostListRouter) Home() {
	this.Data["IsHomePage"] = true
	this.TplNames = "post/home.html"

	//get posts by updated datetime desc order
	var posts []models.Post
	qs := models.Posts()
	cnt, _ := models.CountObjects(qs)
	pager := this.SetPaginator(setting.PostCountPerPage, cnt)
	qs = qs.OrderBy("-Updated").Limit(setting.PostCountPerPage, pager.Offset()).RelatedSel()

	models.ListObjects(qs, &posts)
	this.Data["Posts"] = posts

	//top nav bar data
	var cats []models.Category
	this.setCategories(&cats)
	this.Data["CategorySlug"] = "home"

	//set cookie
	this.Ctx.SetCookie("category_slug", "home", 1<<31-1, "/")
}

//Get the posts by category
func (this *PostListRouter) Category() {
	this.Data["IsHomePage"] = true
	this.Data["IsCategory"] = true
	this.TplNames = "post/home.html"

	//check category slug
	slug := this.GetString(":slug")
	cat := models.Category{Slug: slug}
	if err := cat.Read("Slug"); err != nil {
		this.Abort("404")
		return
	}
	//set cookie
	this.Ctx.SetCookie("category_slug", cat.Slug, 1<<31-1, "/")
	//get posts by category slug, order by Updated desc
	qs := models.Posts().Filter("Category", &cat)
	cnt, _ := models.CountObjects(qs)
	pager := this.SetPaginator(setting.PostCountPerPage, cnt)
	qs = qs.OrderBy("-Updated").Limit(setting.PostCountPerPage, pager.Offset()).RelatedSel()
	var posts []models.Post
	models.ListObjects(qs, &posts)

	this.Data["Category"] = &cat
	this.Data["Posts"] = posts

	//top nav bar data
	var cats []models.Category
	this.setCategories(&cats)
	var topics []models.Topic
	this.setTopicsOfCategory(&topics, &cat)
	this.Data["CategorySlug"] = cat.Slug
}

//Topic Home Page
func (this *PostListRouter) Topic() {
	this.Data["IsHomePage"] = true
	this.Data["IsCategory"] = true
	this.TplNames = "post/topic.html"
	//check topic slug
	slug := this.GetString(":slug")
	topic := models.Topic{Slug: slug}
	if err := topic.Read("Slug"); err != nil {
		this.Abort("404")
		return
	}
	//get topic category
	category := models.Category{Id: topic.Category.Id}
	if err := category.Read("Id"); err != nil {
		this.Abort("404")
		return
	}

	//get posts by topic
	qs := models.Posts().Filter("Topic", &topic)
	cnt, _ := models.CountObjects(qs)
	pager := this.SetPaginator(setting.PostCountPerPage, cnt)
	qs = qs.OrderBy("-Updated").Limit(setting.PostCountPerPage, pager.Offset()).RelatedSel()
	var posts []models.Post
	models.ListObjects(qs, &posts)

	this.Data["Posts"] = posts
	this.Data["Topic"] = &topic
	this.Data["Category"] = &category

	//check whether added it into favorite list
	HasFavorite := false
	if this.IsLogin {
		HasFavorite = models.FollowTopics().Filter("User", &this.User).Filter("Topic", &topic).Exist()
	}
	this.Data["HasFavorite"] = HasFavorite

}

// Add this topic into favorite list
func (this *PostListRouter) TopicSubmit() {
	this.Data["IsHomePage"] = true
	slug := this.GetString(":slug")

	topic := models.Topic{Slug: slug}
	if err := topic.Read("Slug"); err != nil {
		this.Abort("404")
		return
	}

	result := map[string]interface{}{
		"success": false,
	}

	if this.IsAjax() {
		action := this.GetString("action")
		switch action {
		case "favorite":
			if this.IsLogin {
				qs := models.FollowTopics().Filter("User", &this.User).Filter("Topic", &topic)
				if qs.Exist() {
					qs.Delete()
				} else {
					fav := models.FollowTopic{User: &this.User, Topic: &topic}
					fav.Insert()
				}
				topic.RefreshFollowers()
				this.User.RefreshFavTopics()
				result["success"] = true
			}
		}
	}

	this.Data["json"] = result
	this.ServeJson()
}

// Post Router
type PostRouter struct {
	base.BaseRouter
}

func (this *PostRouter) NewPost() {
	this.Data["IsHomePage"] = true
	this.TplNames = "post/new.html"

	if this.CheckActiveRedirect() {
		return
	}

	form := post.PostForm{Locale: this.Locale}
	slug := this.GetString("topic")
	if len(slug) > 0 {
		topic := models.Topic{Slug: slug}
		topic.Read("Slug")
		form.Topic = topic.Id
		form.Category = topic.Category.Id
		this.Data["Topic"] = &topic
	}

	post.ListTopics(&form.Topics)
	this.SetFormSets(&form)
}

func (this *PostRouter) NewPostSubmit() {
	this.Data["IsHomePage"] = true
	this.TplNames = "post/new.html"

	if this.CheckActiveRedirect() {
		return
	}

	form := post.PostForm{Locale: this.Locale}
	slug := this.GetString("topic")
	if len(slug) > 0 {
		topic := models.Topic{Slug: slug}
		topic.Read("Slug")
		form.Topic = topic.Id
		form.Category = topic.Category.Id
		this.Data["Topic"] = &topic
	}

	post.ListTopics(&form.Topics)
	if !this.ValidFormSets(&form) {
		return
	}

	var post models.Post
	if err := form.SavePost(&post, &this.User); err == nil {
		this.JsStorage("deleteKey", "post/new")
		this.Redirect(post.Link(), 302)
	}
}

func (this *PostRouter) loadPost(post *models.Post, user *models.User) bool {
	id, _ := this.GetInt(":post")
	if id > 0 {
		qs := models.Posts().Filter("Id", id)
		if user != nil {
			qs = qs.Filter("User", user.Id)
		}
		qs.RelatedSel(1).One(post)
	}

	if post.Id == 0 {
		this.Abort("404")
		return true
	}

	this.Data["Post"] = post

	return false
}

func (this *PostRouter) loadComments(post *models.Post, comments *[]*models.Comment) {
	qs := post.Comments()
	if num, err := qs.RelatedSel("User").OrderBy("Id").All(comments); err == nil {
		this.Data["Comments"] = *comments
		this.Data["CommentsNum"] = num
	}
}

//Post Page
func (this *PostRouter) SinglePost() {
	this.Data["IsHomePage"] = true
	this.TplNames = "post/post.html"

	var postMd models.Post
	if this.loadPost(&postMd, nil) {
		return
	}

	var comments []*models.Comment
	this.loadComments(&postMd, &comments)

	//mark all notification as read
	if this.IsLogin {
		models.MarkNortificationAsRead(this.User.Id, postMd.Id)
	}

	//check whether this post is favorited
	num, _ := this.User.FavoritePosts().Filter("Post__Id", postMd.Id).Filter("IsFav", true).Count()
	if num != 0 {
		this.Data["IsPostFav"] = true
	} else {
		this.Data["IsPostFav"] = false
	}

	form := post.CommentForm{}
	this.SetFormSets(&form)
	//increment PageViewCount
	post.PostBrowsersAdd(this.User.Id, this.Ctx.Input.IP(), &postMd)
}

//New Comment
func (this *PostRouter) SinglePostCommentSubmit() {
	this.Data["IsHomePage"] = true
	this.TplNames = "post/post.html"

	if this.CheckActiveRedirect() {
		return
	}

	var postMd models.Post
	if this.loadPost(&postMd, nil) {
		return
	}

	var redir bool

	defer func() {
		if !redir {
			var comments []*models.Comment
			this.loadComments(&postMd, &comments)
		}
	}()

	form := post.CommentForm{}
	if !this.ValidFormSets(&form) {
		return
	}

	comment := models.Comment{}
	if err := form.SaveComment(&comment, &this.User, &postMd); err == nil {
		post.FilterCommentMentions(&this.User, &postMd, &comment)
		this.JsStorage("deleteKey", "post/comment")
		this.Redirect(postMd.Link(), 302)
		redir = true

		post.PostReplysCount(&postMd)
	}
}

func (this *PostRouter) EditPost() {
	this.Data["IsHomePage"] = true
	this.TplNames = "post/edit.html"

	if this.CheckActiveRedirect() {
		return
	}

	var postMd models.Post
	if this.loadPost(&postMd, &this.User) {
		return
	}

	form := post.PostForm{}
	form.SetFromPost(&postMd)
	post.ListTopics(&form.Topics)
	this.SetFormSets(&form)
}

func (this *PostRouter) EditPostSubmit() {
	this.Data["IsHomePage"] = true
	this.TplNames = "post/edit.html"

	if this.CheckActiveRedirect() {
		return
	}

	var postMd models.Post
	if this.loadPost(&postMd, &this.User) {
		return
	}

	form := post.PostForm{}
	form.SetFromPost(&postMd)
	post.ListTopics(&form.Topics)
	if !this.ValidFormSets(&form) {
		return
	}

	if err := form.UpdatePost(&postMd, &this.User); err == nil {
		this.JsStorage("deleteKey", "post/edit")
		this.Redirect(postMd.Link(), 302)
	}
}
