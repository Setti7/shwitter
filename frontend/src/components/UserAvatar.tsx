import { Avatar } from "@mui/material";
import { FC } from "react";
import User from "../models/user";
import { Link as RouterLink } from "react-router-dom";

interface Props {
  user: User;
  size?: number;
}

const UserAvatar: FC<Props> = ({ user, size = 36 }) => {
  return (
    <RouterLink to={"/user/" + user.id}>
      <Avatar alt={user.name} sx={{ width: size, height: size }} />
    </RouterLink>
  );
};

export default UserAvatar;
