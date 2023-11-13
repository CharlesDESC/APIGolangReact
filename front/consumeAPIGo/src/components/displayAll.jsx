import React, { useEffect, useState } from "react";
import { Article } from "./articles.jsx";
import "../styles/displayAll.css";

export const DisplayAll = () => {
  const [rows, setRows] = useState([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);

  const pageLimit = (page) => {
    if (page < 1) {
      return 1;
    }
    return page;
  };

  const fetchData = async () => {
    try {
      const response = await fetch(
        `http://localhost:8080/v1/articles?page=${pageLimit(page)}`
      );
      if (!response.ok) {
        throw new Error(`Network response was not ok: ${response.statusText}`);
      }
      const json = await response.json();
      setRows(json || []);
    } catch (error) {
      console.error("Error fetching data:", error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData(); // Appeler la fonction fetchData immédiatement après le montage du composant
  }, [page]);

  if (loading) {
    return <p>Loading...</p>; // Vous pouvez remplacer cela par un indicateur de chargement plus sophistiqué si nécessaire
  }

  return (
    <div className="display-all">
      <button onClick={() => setPage(page - 1)}>Previous</button>
      <div className="articles">
        {rows.map((row) => (
          <Article key={row.id} {...row} />
        ))}
      </div>
      <button onClick={() => setPage(page + 1)}>Next</button>
    </div>
  );
};
