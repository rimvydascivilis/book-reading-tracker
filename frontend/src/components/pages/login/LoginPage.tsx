import React from "react";
import {Navigate} from "react-router-dom";
import {
  GoogleOAuthProvider,
  GoogleLogin,
  CredentialResponse,
} from "@react-oauth/google";
import api from "../../../api/api";
import PathConstants from "../../../routes/PathConstants";
import config from "../../../config/config";
import {message} from "antd";
import {useAuth} from "../../../context/AuthContext";
import {IAxiosError} from "../../../types/errorTypes";

const LoginPage: React.FC = () => {
  const {login, isAuthenticated} = useAuth();

  const handleLoginSuccess = async (credentialResponse: CredentialResponse) => {
    if (!credentialResponse.credential) {
      message.error("Failed to login: No credential received.");
      return;
    }

    try {
      const response = await api.post("/auth/login", {
        token: credentialResponse.credential,
      });
      login(response.data.token);
      message.success("Login successful!");
    } catch (error) {
      const axiosError = error as IAxiosError;
      message.error(
        "Failed to login: " +
          (axiosError.response?.data.message || "Network error"),
      );
    }
  };

  if (isAuthenticated) {
    return <Navigate to={PathConstants.HOME} />;
  }

  return (
    <div style={{display: "flex", justifyContent: "center", height: "100%"}}>
      <GoogleOAuthProvider clientId={config.googleClientId}>
        <div style={{textAlign: "center", maxWidth: 400}}>
          <h1>Login with Google</h1>
          <GoogleLogin onSuccess={handleLoginSuccess} context="signin" />
        </div>
      </GoogleOAuthProvider>
    </div>
  );
};

export default LoginPage;
