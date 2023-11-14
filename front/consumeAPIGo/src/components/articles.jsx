import "../styles/articles.css";

export const Article = ({ title, content, author, id }) => {
  const modify = () => {
    console.log("Modify article with ID:", id);
  };

  const deleteArticle = () => {
    console.log("Delete article with ID:", id);
    fetch(`http://127.0.0.1:8080/v1/articles/${id}`, {
      method: "DELETE",
    })
      .then((response) => {
        if (response.ok) {
          console.log("Article deleted successfully");
          alert("Article deleted successfully");
        } else {
          console.error("Error deleting article");
          alert("Error deleting article");
        }
      })
      .catch((error) => {
        console.error("Network error while deleting article", error);
      });
  };

  return (
    <div className="">
      <h1>titre : {title}</h1>
      <p>content : {content}</p>
      <p>auteur : {author}</p>
      <button className="button1" onClick={() => modify()}>
        Modifier
      </button>
      <button className="button3" onClick={() => deleteArticle()}>
        Supprimer
      </button>
    </div>
  );
};
