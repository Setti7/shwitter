export default interface Session {
  token: string;
}

export const getSessionID = (session: Session) => {
  return session.token.split(":").pop();
};
