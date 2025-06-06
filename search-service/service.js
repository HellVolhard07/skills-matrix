import { users } from './data.js';

export function searchUsersBySkill(call, callback) {
  const { skill_name, category, proficiency } = call.request;

  const results = users.filter(user =>
    user.skills.some(skill =>
      (!skill_name || skill.name.toLowerCase().includes(skill_name.toLowerCase())) &&
      (!category || skill.category.toLowerCase() === category.toLowerCase()) &&
      (!proficiency || skill.proficiency.toLowerCase() === proficiency.toLowerCase())
    )
  );

  const userResults = results.map(user => ({
    id: user.id,
    name: user.name,
    department: user.department,
    skills: user.skills.map(s => s.name)
  }));

  callback(null, { users: userResults });
}
