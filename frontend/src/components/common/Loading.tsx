import React from "react";
import {Spin} from "antd";

const Loading: React.FC = () => {
  return (
    <div
      style={{
        width: "100%",
        height: "100%",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
      }}>
      <Spin size="large" />
    </div>
  );
};

export default Loading;
