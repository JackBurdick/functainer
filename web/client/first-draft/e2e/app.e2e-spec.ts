import { FirstDraftPage } from './app.po';

describe('first-draft App', () => {
  let page: FirstDraftPage;

  beforeEach(() => {
    page = new FirstDraftPage();
  });

  it('should display welcome message', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('Welcome to app!');
  });
});
