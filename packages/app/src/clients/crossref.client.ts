import { tryCatch } from '@doi/utils/try-catch';

type Author = { given: string; family: string };

type CrossReferenceResponse = {
  status: string;
  message: {
    author: { given: string; family: string }[];
    title: string[];
    'container-title': string[];
    volume: string;
    issue: string;
    page: string;
    'published-online': { 'date-parts': number[][] };
    'published-print': { 'date-parts': number[][] };
  };
};

export type Reference = {
  id: string;
  authors: Author[];
  title: string;
  journal: string;
  volume: string;
  issue: string;
  pages: string;
  year: number;
  url: string;
};

export const getId = (url: string): string | null => {
  try {
    // Normalize and remove prefix
    return url.replace(/^https?:\/\/(dx\.)?doi\.org\//, '');
  } catch {
    return null;
  }
};

export const getWork = async (
  id: string,
  {
    timeoutMs = 60000,
  }: {
    timeoutMs?: number;
  } = { timeoutMs: 60000 },
): Promise<{ reference: Reference | null }> => {
  const abortController: AbortController = new AbortController();
  const timeout = setTimeout(() => abortController.abort(), timeoutMs);
  const url: string = `https://api.crossref.org/works/${id}`;
  const urlSearchParams: URLSearchParams = new URLSearchParams();
  urlSearchParams.append('url', url);
  const proxyUrl: string = `https://hieudoanm-proxy.vercel.app/api/?${urlSearchParams.toString()}`;
  console.info('proxyUrl', proxyUrl);
  const method: string = 'GET';
  const headers: HeadersInit = { 'Content-Type': 'application/json' };
  const { data: response, error } = await tryCatch(
    fetch(proxyUrl, { method, headers, signal: abortController.signal }),
  );
  if (error) {
    clearTimeout(timeout);
    console.error('Error in generateContent:', error);
    return { reference: null };
  }
  if (!response?.ok) {
    clearTimeout(timeout);
    console.error(
      `Failed to generate content: ${response?.statusText ?? 'Unknown error'}`,
    );
    return { reference: null };
  }
  const data: CrossReferenceResponse = await response.json();
  clearTimeout(timeout);
  const { message } = data;
  const published = message['published-print'] ??
    message['published-online'] ?? { 'date-parts': [[]] };
  return {
    reference: {
      id,
      authors: message.author ?? [],
      title: message.title.at(0) ?? '',
      journal: message['container-title'].at(0) ?? '',
      volume: message.volume,
      issue: message.issue,
      pages: message.page,
      year: published['date-parts'].at(0)?.at(0) ?? 0,
      url: `https://doi.org/${id}`,
    },
  };
};
