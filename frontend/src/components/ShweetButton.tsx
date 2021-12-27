import { Fab } from "@mui/material";
import HistoryEduIcon from "@mui/icons-material/HistoryEdu";
import { useNavigate } from "react-router-dom";
import { NEW_SHWEET_ROUTE } from "../config/routes";

const ShweetButton = () => {
  const navigate = useNavigate();

  return (
    <Fab
      sx={{
        position: "absolute",
        bottom: 16,
        right: 16,
      }}
      variant="extended"
      onClick={() => navigate(NEW_SHWEET_ROUTE)}
    >
      <HistoryEduIcon sx={{ mr: 1 }} />
      Shweet
    </Fab>
  );
};

export default ShweetButton;
