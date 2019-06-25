let chai = require("chai");
let expect = chai.expect;
chai.should();

let chaiHttp = require("chai-http");
chai.use(chaiHttp);

let server = "http://localhost:9090";

it("it should be running", done => {
  chai
    .request(server)
    .get("/metrics")
    .end((err, res) => {
      expect(err).to.be.null;
      res.should.have.status(200);
      res.should.have.header("content-type", /^text\/plain/);
      done();
    });
});

it("it should contain pipepine information", done => {
  chai
    .request(server)
    .get("/metrics")
    .end((err, res) => {
      expect(err).to.be.null;
      res.text.should.contain("gocd_pipeline");
      done();
    });
});
