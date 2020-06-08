import { makeStyles } from "@material-ui/core/styles";

const useStyles = makeStyles({
  copyright: {
    // width: "50%",
    left: "50%",
    position: "relative",
    top: "50%",
    marginRight: "-50%",
    transform: "translate(-50%, -50%)",
    textAlign: "center",
  },
  footer: {
    height: "calc(10vh + max-content)",
    display: "inline",
    // paddingTop: "5vh",
  },
});
export default function Footer() {
  const classes = useStyles();
  return (
    <div className={classes.footer}>
      <div className={classes.copyright}>
        <p>
          &copy;
          {` ${new Date().getFullYear()}`}
          {", "}
          Gordon Pham-Nguyen
        </p>
        <p>
          {"Source code on "}
          <a
            href="https://github.com/gordonpn/internet-speedtests-visualized"
            target="_blank"
            rel="noopener noreferrer"
          >
            GitHub
          </a>
        </p>
        <p>
          {"More about me on my "}
          <a
            href="https://gordon-pn.com/"
            target="_blank"
            rel="noopener noreferrer"
          >
            Website
          </a>
        </p>
      </div>
    </div>
  );
}
