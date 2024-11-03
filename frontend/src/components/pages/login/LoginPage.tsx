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
import {useAuth} from "../../../context/AuthContext";

const LoginPage: React.FC = () => {
  const {login, isAuthenticated} = useAuth();

  const handleLoginSuccess = (credentialResponse: CredentialResponse) => {
    if (!credentialResponse.credential) {
      console.error("No credential response");
      return;
    }
    api
      .post("/auth/login", {token: credentialResponse.credential})
      .then(response => {
        if (response.status === 200) {
          login(response.data.token);
          window.location.reload();
        }
      });
  };

  if (isAuthenticated) {
    return <Navigate to={PathConstants.LIBRARY} />;
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
