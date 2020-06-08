import { makeStyles } from "@material-ui/core/styles";
import Paper from "@material-ui/core/Paper";

const useStyles = makeStyles({
  header: {
    background: "#3b4a6b",
    backgroundPosition: "center",
    backgroundRepeat: "no-repeat",
    backgroundSize: "cover",
    height: "15vh",
  },
  title: {
    width: "50%",
    left: "50%",
    position: "relative",
    top: "50%",
    marginRight: "-50%",
    transform: "translate(-50%, -50%)",
    textAlign: "center",
    color: "#ffb6c1",
  },
});

export default function Header() {
  const classes = useStyles();
  return (
    <Paper elevation={24} className={classes.header}>
      <h1 className={classes.title}>Hot Flag Deals</h1>
    </Paper>
  );
}
