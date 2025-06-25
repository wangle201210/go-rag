import hljs from 'highlight.js';
import { marked } from 'marked';
import 'highlight.js/styles/github.css';

// Initialize marked configuration
marked.setOptions({
  highlight(code, lang) {
    const language = hljs.getLanguage(lang) ? lang : 'plaintext';
    return hljs.highlight(code, { language }).value;
  },
  langPrefix: 'hljs language-',
  gfm: true,
  breaks: true,
});

export function renderMarkdown(text) {
  if (!text) {
    return '';
  }
  return marked.parse(text);
}