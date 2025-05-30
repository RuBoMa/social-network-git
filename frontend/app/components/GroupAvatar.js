import React from 'react';
import Link from 'next/link';

export default function GroupAvatar({ group, size = "md", disableLink = false, truncateName = false }) {
  if (!group) return null;

  const sizeClasses = {
    xs: "w-4 h-4 text-xs",
    sm: "w-6 h-6 text-sm",
    md: "w-8 h-8 text-base",
    lg: "w-10 h-10 text-lg",
    xl: "w-16 h-16 text-xl",
  };

  const imageUrl = '/group-default.png';
  const displayName = group.group_name
    ? (truncateName && group.group_name.length > 18
        ? group.group_name.slice(0, 18) + "..."
        : group.group_name)
    : 'Unnamed Group';

  if (disableLink) {
    return (
      <div className="flex items-center gap-2">
        <img
          src={imageUrl}
          alt="Group avatar"
          className={`rounded-full ${sizeClasses[size] || sizeClasses.md}`}
        />
        <span className="font-medium">{displayName}</span>
      </div>
    )
  }

  return (
    <div className="flex items-center gap-2">
      <img
        src={imageUrl}
        alt="Group avatar"
        className={`rounded-full ${sizeClasses[size] || sizeClasses.md}`}
      />
      <Link
        href={`/group?group_id=${group.group_id}`}>
         <span className="font-medium hover:underline">{displayName}</span>
      </Link>
    </div>
  )
}