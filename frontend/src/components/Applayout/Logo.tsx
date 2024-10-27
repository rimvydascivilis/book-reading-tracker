import React from "react";
import logo from "../../assets/logo.png";
import {Link} from "react-router-dom";

const Header: React.FC = () => {
  return (
    <Link to="/">
      <div
        style={{
          display: "flex",
          alignItems: "center",
          background: "white",
          padding: "0",
        }}>
        <img
          src={logo}
          alt="logo"
          style={{width: "50px", height: "50px", margin: "0 15px"}}
        />
        <span style={{fontSize: "24px", fontWeight: "bold"}}>Book tracker</span>
      </div>
    </Link>
  );
};

export default Header;
