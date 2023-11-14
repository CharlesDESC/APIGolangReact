import React, { useEffect, useState } from "react";
import { Article } from "./articles.jsx";
import "../styles/displayAll.css";

export const DisplayAll = () => {
  const [rows, setRows] = React.useState([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);

  const pageLimit = (page) => {
    if (page < 1) {
      setPage(1);
    }
    return page;
  };

  useEffect(() => {
    setLoading(true);
    fetch(`http://localhost:8080/v1/articles?page=${pageLimit(page)}`)
      .then((data) => data.json())
      .then((json) => setRows(json || []))
      .then(() => setLoading(false));
  }, [page]);

  console.log("rows", rows);
  if (loading) {
    return <p>Loading...</p>;
  }

  return (
    <div className="display-all">
      <button className="button1" onClick={() => setPage(page - 1)}>
        Previous
      </button>
      <div className="articles">
        {rows.articles.map((row, key) => {
          return (
            <div className="article" key={key}>
              <Article
                title={row.title}
                content={row.content}
                author={row.author}
                id={row.article_id}
              />
            </div>
          );
        })}
      </div>
      <button className="button2" onClick={() => setPage(page + 1)}>
        Next
      </button>
    </div>
  );
};
