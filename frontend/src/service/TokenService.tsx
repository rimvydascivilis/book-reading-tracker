const tokenKey = "token";

export function setToken(token: string): void {
  localStorage.setItem(tokenKey, token);
}

export const removeToken = (): void => {
  localStorage.removeItem(tokenKey);
};

export const getToken = (): string | null => {
  return localStorage.getItem(tokenKey);
};

export const isTokenSet = (): boolean => {
  return localStorage.getItem(tokenKey) !== null;
};
