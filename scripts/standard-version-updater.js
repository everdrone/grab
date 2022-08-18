module.exports.readVersion = (contents) => {
  const re = /const\sVersion\s=\s"(.*)"/;

  if (!contents.match(re)[1]) {
    throw new Error("Could not find version.");
  }

  console.log("--------------");
  console.log(contents.match(re)[1]);

  return contents.match(re)[1];
};

module.exports.writeVersion = (contents, version) => {
  const re = /const\sVersion\s=\s"(.*)"/;
  const oldVersion = contents.match(re)[1];

  console.log("--------------");
  console.log(oldVersion, version);
  console.log(contents.replace(oldVersion, version));

  return contents.replace(oldVersion, version);
};
