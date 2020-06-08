import Paper from "@material-ui/core/Paper";
import { makeStyles } from "@material-ui/core/styles";
import Footer from "./Footer";

const useStyles = makeStyles({
  dealsBody: {
    background: "#fff",
    height: "100vh",
  },
});
export default function Content() {
  const classes = useStyles();
  return <Paper elevation={24} className={classes.dealsBody} />;
}
