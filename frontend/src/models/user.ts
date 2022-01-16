export default interface User {
  id: string;
  username: string;
  name: string;
  email: string;
  bio: string;
}

export interface UserProfile {
  id: string;
  username: string;
  name: string;
  email: string;
  bio: string;
  followers_count: number;
  friends_count: number;
  shweets_count: number;
}
