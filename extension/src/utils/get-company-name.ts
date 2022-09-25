export default async function getCompanyName(domain: string): Promise<string | undefined> {
  const text = await (
    await fetch(`https://www.similarweb.com/website/${domain}`)
  ).text();
  return JSON.parse(/window.__APP_DATA__ = (.*)$/m.exec(text)[1]).layout.data.overview.companyName;
}
