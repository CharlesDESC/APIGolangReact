import React, { useState } from "react";
import "../styles/addNewArticle.css";

export const AddNewArticle = () => {
  const [visible, setVisible] = useState(false);
  const [article, setArticle] = useState({
    title: "",
    slug: "",
    content: "",
    author: "",
    link: "",
  });

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setArticle((prevArticle) => ({
      ...prevArticle,
      [name]: value,
    }));
  };

  const handlePostArticle = () => {
    console.log("article", article);
    fetch("http://127.0.0.1:8080/v1/articles", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        title: article.title,
        slug: article.slug,
        content: article.content,
      }),
    })
      .then((response) => response.json())
      .then((data) => {
        console.log("Success:", data);
        // Handle the response data as needed
      })
      .catch((error) => {
        console.error("Error:", error);
        // Handle errors
      });
  };

  return (
    <div>
      <button className="new" onClick={() => setVisible(!visible)}>
        <h1>Add New Article</h1>
      </button>

      {visible && (
        <div>
          <form>
            <input
              type="text"
              placeholder="Title"
              name="title"
              value={article.title}
              onChange={handleInputChange}
            />
            <input
              type="text"
              placeholder="Content"
              name="content"
              value={article.content}
              onChange={handleInputChange}
            />
            <input
              type="text"
              placeholder="Author"
              name="author"
              value={article.author}
              onChange={handleInputChange}
            />
            <button
              className="button1"
              type="button"
              onClick={handlePostArticle}
            >
              Submit
            </button>
          </form>
        </div>
      )}
    </div>
  );
};
