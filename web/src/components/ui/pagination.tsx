import { ChevronLeft, ChevronRight } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";

interface PaginationProps {
  page: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

function pageItems(page: number, totalPages: number) {
  const pages = new Set([1, totalPages, page - 1, page, page + 1]);
  return Array.from(pages)
    .filter((p) => p >= 1 && p <= totalPages)
    .sort((a, b) => a - b);
}

export function Pagination({
  page,
  totalPages,
  onPageChange,
}: PaginationProps) {
  const pages = pageItems(page, totalPages);

  return (
    <nav
      className="flex items-center justify-center gap-1"
      aria-label="Pagination"
    >
      <Button
        type="button"
        variant="outline"
        size="icon"
        className="h-8 w-8"
        disabled={page <= 1}
        onClick={() => onPageChange(page - 1)}
        aria-label="Previous page"
      >
        <ChevronLeft className="h-4 w-4" />
      </Button>

      {pages.map((item, index) => {
        const previous = pages[index - 1];
        return (
          <div key={item} className="flex items-center gap-1">
            {previous && item - previous > 1 && (
              <span className="flex h-8 w-8 items-center justify-center text-xs text-muted-foreground">
                ...
              </span>
            )}
            <Button
              type="button"
              variant={item === page ? "default" : "outline"}
              size="icon"
              className={cn(
                "h-8 w-8 text-xs",
                item === page && "pointer-events-none",
              )}
              onClick={() => onPageChange(item)}
              aria-current={item === page ? "page" : undefined}
            >
              {item}
            </Button>
          </div>
        );
      })}

      <Button
        type="button"
        variant="outline"
        size="icon"
        className="h-8 w-8"
        disabled={page >= totalPages}
        onClick={() => onPageChange(page + 1)}
        aria-label="Next page"
      >
        <ChevronRight className="h-4 w-4" />
      </Button>
    </nav>
  );
}
