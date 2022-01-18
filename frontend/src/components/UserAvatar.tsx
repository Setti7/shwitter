import { Avatar } from "@mui/material";
import { FC } from "react";
import User from "../models/user";
import { Link } from "react-router-dom";

interface Props {
  user: User;
  size?: number;
}

const UserAvatar: FC<Props> = ({ user, size = 36 }) => {
  return (
    <Link
      to={"/user/" + user.id}
      style={{ textDecoration: "none", color: "white" }}
    >
      <Avatar
        alt={user.name}
        sx={{ width: size, height: size }}
        src="/broken.png"
      />
    </Link>
  );
};

export default UserAvatar;
