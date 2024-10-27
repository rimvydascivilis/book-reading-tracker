import {Button} from "antd";
import React from "react";
import api from "../../../api/api";

const HomePage: React.FC = () => {
  const handleClick = () => {
    console.log("Button clicked");
    api.post("/").then(response => {
      console.log(response);
    });
  };
  return (
    <>
      <h1>Home</h1>
      <Button type="primary" onClick={handleClick}>
        Primary Button
      </Button>
    </>
  );
};

export default HomePage;
